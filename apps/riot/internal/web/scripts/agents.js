// agents browser client: lists every agent alphabetically, filters them from the
// search needle, and opens one into the main view via the ?agent= search param.
(function () {
	"use strict";

	// the narrowest the list column may be dragged, and the widest it may grow
	// as a fraction of the viewport.
	var MIN_WIDTH = 280;
	var MAX_FRACTION = 0.6;

	var state = {
		agents: [],
		filter: "",
		selected: "",
	};

	var query = async (document, variables) => {
		try {
			var response = await fetch("/graph", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ query: document, variables: variables || {} }),
			});
			var payload = await response.json();
			if (!response.ok) return { data: null, error: payload.errors || payload };
			if (payload.errors && payload.errors.length) return { data: payload.data, error: payload.errors };
			return { data: payload.data, error: null };
		} catch (err) {
			return { data: null, error: err };
		}
	};

	// DEPENDENCY_DEPTH is how many levels of dependencies the tree query asks
	// for. The catalogue nests about five deep, so this is comfortably more than
	// any real chain; a skill with no dependencies returns an empty array, which
	// ends the recursion in the data.
	var DEPENDENCY_DEPTH = 30;

	// dependencyFragment returns a GraphQL fragment selecting a skill and its
	// dependencies to the given depth. Each level nests the one before it, so a
	// single query returns the whole chain rather than one level of it.
	var dependencyFragment = function (depth) {
		var fields = "id name description";
		for (var i = 0; i < depth; i++) {
			fields = "id name description dependencies { " + fields + " }";
		}
		return "fragment SkillTree on Skill { " + fields + " }";
	};

	var loadAgents = async () => {
		var result = await query(
			"{ Agents { id name description skills { ...SkillTree } } } " + dependencyFragment(DEPENDENCY_DEPTH)
		);
		if (result.error || !result.data) {
			console.error("Agents: failed getting data", result.error);
			return null;
		}
		return result.data.Agents || [];
	};

	// byName orders the list alphabetically, case-insensitively, so the column
	// reads the same however the catalogue is stored.
	var byName = function (a, b) {
		return String(a.name || "").localeCompare(String(b.name || ""), undefined, { sensitivity: "base" });
	};

	var renderList = function () {
		var list = document.querySelector("#agents-items");
		if (!list) return;

		var term = state.filter.trim().toLowerCase();
		var agents = state.agents
			.filter(function (agent) {
				if (!term) return true;
				return String(agent.name || "").toLowerCase().includes(term);
			})
			.sort(byName);

		list.textContent = "";

		if (agents.length === 0) {
			var empty = document.createElement("div");
			empty.className = "agents-empty";
			empty.textContent = term ? "no matching agents" : "no agents";
			list.appendChild(empty);
			return;
		}

		agents.forEach(function (agent) {
			var button = document.createElement("button");
			button.type = "button";
			button.className = "agent-item";
			button.dataset.agent = agent.id;
			button.textContent = agent.name;
			button.title = agent.description || agent.name;
			button.classList.toggle("active", agent.id === state.selected);
			button.addEventListener("click", function () {
				selectAgent(agent.id);
			});
			list.appendChild(button);
		});
	};

	// selectAgent opens an agent in the main view and records it in the url, so
	// the current agent survives a reload and can be linked to.
	var selectAgent = function (id, options) {
		var push = !options || options.push !== false;
		var agent = state.agents.find(function (candidate) {
			return candidate.id === id;
		});
		if (!agent) return;

		state.selected = agent.id;

		var heading = document.querySelector("#agent-name");
		if (heading) heading.textContent = agent.name;

		renderList();
		renderAgentTree(agent);

		var url = new URL(window.location.href);
		url.searchParams.set("agent", agent.id);
		if (push) window.history.pushState({ agent: agent.id }, "", url.toString());
		else window.history.replaceState({ agent: agent.id }, "", url.toString());
	};

	// buildAgentTree turns an agent into a hierarchy rooted at the agent itself:
	// its skills are the first branch, and each skill's dependencies hang beneath
	// it. A skill is placed at most once, so a dependency shared by several
	// skills appears under the first one that claims it and a cycle cannot
	// recurse forever.
	var buildAgentTree = function (agent) {
		var skills = agent.skills || [];

		// index every skill we can see, including ones only reachable as a
		// dependency, so a nested stub can be replaced by its full record.
		var byId = new Map();
		var index = function (skill) {
			if (!skill || !skill.id) return;
			if (!byId.has(skill.id)) byId.set(skill.id, skill);
			(skill.dependencies || []).forEach(index);
		};
		skills.forEach(index);

		var placed = new Set();
		var build = function (skill) {
			placed.add(skill.id);
			var node = {
				id: skill.id,
				name: skill.name,
				description: skill.description || "",
				children: [],
			};
			(skill.dependencies || []).forEach(function (dep) {
				if (placed.has(dep.id)) return;
				node.children.push(build(byId.get(dep.id) || dep));
			});
			node.children.sort(byName);
			return node;
		};

		var root = { id: agent.id, name: agent.name, description: agent.description || "", children: [] };
		skills.slice().sort(byName).forEach(function (skill) {
			if (placed.has(skill.id)) return;
			root.children.push(build(skill));
		});
		return root;
	};

	// renderAgentTree draws the agent's skills as an indented left-to-right tree,
	// the same layout the skills page's "tree" view uses: the agent's name sits
	// at the top left and the branches fan out beneath and to the right.
	var renderAgentTree = function (agent) {
		var stage = document.querySelector("#agent-graph");
		if (!stage) return;

		stage.textContent = "";

		if (typeof d3 === "undefined") {
			var missing = document.createElement("p");
			missing.className = "errored";
			missing.textContent = "d3 failed to load";
			stage.appendChild(missing);
			return;
		}

		var tree = buildAgentTree(agent);
		if (tree.children.length === 0) {
			var empty = document.createElement("p");
			empty.className = "agents-empty";
			empty.textContent = "this agent has no skills";
			stage.appendChild(empty);
			return;
		}

		var fontSize = 16;
		var rowHeight = 32;
		var columnWidth = 170;

		// a label is truncated to what fits between its own node and the next
		// column, so a deep chain never writes over the column to its right.
		var labelChars = Math.max(Math.floor((columnWidth - 16) / (fontSize * 0.55)), 8);

		var root = d3.hierarchy(tree);

		// rows are assigned in pre-order, so every child sits on a row below its
		// parent. d3.tree() centres a parent among its children, which can place
		// the first child above it; the indented layout must always step down.
		var row = 0;
		root.eachBefore(function (d) {
			d.row = row++;
		});

		var svg = d3.select(stage)
			.append("svg")
			.attr("preserveAspectRatio", "xMinYMin meet")
			.attr("role", "img")
			.attr("aria-label", "tree of " + agent.name + "'s skills and their dependencies");

		// the agent is the root, so it sits in the first column at the top left.
		var posX = function (d) {
			return d.depth * columnWidth;
		};

		var posY = function (d) {
			return d.row * rowHeight;
		};

		// links: a step-before path, horizontal to the child's column then
		// vertical down to its row.
		var link = svg.append("g")
			.attr("class", "links")
			.attr("fill", "none")
			.attr("stroke", "currentColor")
			.attr("stroke-opacity", 0.35)
			.attr("stroke-width", 1.5)
			.selectAll("path")
			.data(root.links())
			.join("path")
			.attr("d", function (d) {
				return "M" + posX(d.source) + "," + posY(d.source) + "H" + posX(d.target) + "V" + posY(d.target);
			});

		var node = svg.append("g")
			.attr("class", "nodes")
			.selectAll("g")
			.data(root.descendants())
			.join("g")
			.attr("class", "node")
			.attr("transform", function (d) {
				return "translate(" + posX(d) + "," + posY(d) + ")";
			});

		node.append("circle")
			.attr("class", "hit")
			.attr("r", 9)
			.attr("fill", "transparent");

		node.append("circle")
			.attr("r", function (d) {
				return d.children ? 3.5 : 2.5;
			})
			.attr("fill", function (d) {
				return d.children ? "currentColor" : "#8a8a8a";
			});

		node.append("text")
			.attr("dy", "0.31em")
			.attr("x", 8)
			.attr("paint-order", "stroke")
			.attr("stroke", "#fff")
			.attr("stroke-width", 3)
			.attr("stroke-linejoin", "round")
			.attr("font-size", fontSize)
			.attr("font-family", "sans-serif")
			.text(function (d) {
				return truncate(d.data.name, labelChars);
			});

		// the root is the agent itself, so it is drawn as the tree's heading.
		node.filter(function (d) {
			return d.depth === 0;
		}).classed("agent-root", true);

		node.append("title")
			.text(function (d) {
				return d.data.description ? d.data.name + "\n\n" + d.data.description : d.data.name;
			});

		// hovering a node highlights the path back to the agent.
		node.on("mouseenter", function (event, d) {
			var path = new Set(d.ancestors());
			link.attr("stroke-opacity", function (l) {
				return path.has(l.target) ? 0.95 : 0.12;
			}).attr("stroke-width", function (l) {
				return path.has(l.target) ? 2.5 : 1.5;
			});
			node.classed("active", function (n) {
				return path.has(n);
			});
		}).on("mouseleave", function () {
			link.attr("stroke-opacity", 0.35).attr("stroke-width", 1.5);
			node.classed("active", false);
		});

		var caption = document.createElement("p");
		caption.className = "caption";
		caption.textContent = (root.descendants().length - 1) + " skills · " + root.leaves().length + " leaves";
		stage.appendChild(caption);

		fitAndZoom(stage, svg);
	};

	// truncate shortens a label so long skill names do not swamp the drawing.
	// The full name is always available in the node's <title>.
	var truncate = function (name, maxChars) {
		var text = String(name || "");
		return text.length > maxChars ? text.slice(0, maxChars - 1) + "…" : text;
	};

	// fitAndZoom fits the viewBox to the drawn content, then wires wheel zoom
	// about the pointer, middle-drag pan and double-click to reset. The home
	// view fits the whole tree, so no branch is ever hidden on load; the wheel
	// zooms in to read a label.
	var fitAndZoom = function (stage, svg) {
		var bounds = svg.node().getBBox();
		var pad = 12;
		var baseView = [bounds.x - pad, bounds.y - pad, bounds.width + pad * 2, bounds.height + pad * 2];

		var home = { k: 1, cx: baseView[0] + baseView[2] / 2, cy: baseView[1] + baseView[3] / 2 };
		var zoom = home;

		var viewFor = function (z) {
			return [z.cx - baseView[2] / (2 * z.k), z.cy - baseView[3] / (2 * z.k), baseView[2] / z.k, baseView[3] / z.k];
		};

		var svgNode = svg.node();
		svg.attr("viewBox", viewFor(zoom).join(" "));

		svgNode.addEventListener("wheel", function (event) {
			event.preventDefault();

			var ctm = svgNode.getScreenCTM();
			if (!ctm) return;
			var point = new DOMPoint(event.clientX, event.clientY).matrixTransform(ctm.inverse());

			var k = Math.min(Math.max(zoom.k * Math.exp(-event.deltaY * 0.0015), 0.4), 12);
			if (k === zoom.k) return;

			var ratio = zoom.k / k;
			zoom = {
				k: k,
				cx: point.x + (zoom.cx - point.x) * ratio,
				cy: point.y + (zoom.cy - point.y) * ratio,
			};
			svg.attr("viewBox", viewFor(zoom).join(" "));
		}, { passive: false });

		svgNode.addEventListener("dblclick", function (event) {
			if (event.target !== svgNode) return;
			zoom = home;
			svg.attr("viewBox", viewFor(zoom).join(" "));
		});

		// middle-mouse drag pans the view, complementing the wheel zoom.
		var panning = false;
		var panStart = null;

		var onPanMove = function (event) {
			if (!panning || !panStart) return;
			event.preventDefault();

			var ctm = svgNode.getScreenCTM();
			if (!ctm) return;
			var inv = ctm.inverse();

			var dx = event.clientX - panStart.x;
			var dy = event.clientY - panStart.y;
			var baseDx = inv.a * dx + inv.c * dy;
			var baseDy = inv.b * dx + inv.d * dy;

			zoom = { k: panStart.zoom.k, cx: panStart.zoom.cx - baseDx, cy: panStart.zoom.cy - baseDy };
			svg.attr("viewBox", viewFor(zoom).join(" "));
		};

		var onPanUp = function (event) {
			if (!panning) return;
			panning = false;
			panStart = null;
			svgNode.classList.remove("panning");
			svgNode.releasePointerCapture?.(event.pointerId);
			window.removeEventListener("pointermove", onPanMove);
			window.removeEventListener("pointerup", onPanUp);
			window.removeEventListener("pointercancel", onPanUp);
		};

		svgNode.addEventListener("pointerdown", function (event) {
			if (event.button !== 1) return;
			event.preventDefault();

			panning = true;
			panStart = { x: event.clientX, y: event.clientY, zoom: zoom };
			svgNode.classList.add("panning");
			svgNode.setPointerCapture?.(event.pointerId);
			window.addEventListener("pointermove", onPanMove);
			window.addEventListener("pointerup", onPanUp);
			window.addEventListener("pointercancel", onPanUp);
		});

		svgNode.addEventListener("auxclick", function (event) {
			if (event.button === 1) event.preventDefault();
		});
		svgNode.addEventListener("mousedown", function (event) {
			if (event.button === 1) event.preventDefault();
		});
	};

	// wireSearch filters the list as the needle is typed. The filter is local,
	// so it never touches the url: only the selected agent is linkable.
	var wireSearch = function () {		var input = document.querySelector("#agents-search");
		if (!input) return;
		input.addEventListener("input", function () {
			state.filter = input.value;
			renderList();
		});
	};

	// wireResizer makes the separator drag the boundary between the list and the
	// main view. Only the list's width is set; the main view takes the rest.
	var wireResizer = function () {
		var handle = document.querySelector("#agents-resizer");
		var list = document.querySelector("#agents-list");
		if (!handle || !list) return;

		var startX = 0;
		var startWidth = 0;

		var onMove = function (event) {
			var delta = event.clientX - startX;
			var max = Math.max(window.innerWidth * MAX_FRACTION, MIN_WIDTH);
			var width = Math.min(Math.max(startWidth + delta, MIN_WIDTH), max);
			list.style.width = width + "px";
		};

		var onUp = function () {
			handle.classList.remove("dragging");
			window.removeEventListener("pointermove", onMove);
			window.removeEventListener("pointerup", onUp);
		};

		handle.addEventListener("pointerdown", function (event) {
			event.preventDefault();
			startX = event.clientX;
			startWidth = list.getBoundingClientRect().width;
			handle.classList.add("dragging");
			window.addEventListener("pointermove", onMove);
			window.addEventListener("pointerup", onUp);
		});
	};

	var main = async function () {
		wireSearch();
		wireResizer();

		var agents = await loadAgents();
		if (!agents) {
			var list = document.querySelector("#agents-items");
			if (list) {
				list.textContent = "";
				var p = document.createElement("p");
				p.className = "errored";
				p.textContent = "failed to load agents";
				list.appendChild(p);
			}
			return;
		}

		state.agents = agents;
		renderList();

		// ?agent=<identifier> opens straight onto that agent.
		var params = new URLSearchParams(window.location.search);
		var wanted = String(params.get("agent") || "").trim();
		if (wanted) selectAgent(wanted, { push: false });

		window.addEventListener("popstate", function () {
			var next = new URLSearchParams(window.location.search);
			var id = String(next.get("agent") || "").trim();
			if (id) selectAgent(id, { push: false });
		});
	};

	document.addEventListener("DOMContentLoaded", function () {
		main();
	});
})();
