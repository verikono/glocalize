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

		// the skill label size, and the agent's name at 150% of it.
		var SKILL_FONT_PX = 12;
		var ROOT_FONT_PX = SKILL_FONT_PX * 1.5;

		// the agent's name starts here, leaving room above and to its left.
		var ORIGIN_X = 100;
		var ORIGIN_Y = 100;

		// the gap between a label's right edge and the next column's node, and
		// the gap between one row's text and the next.
		var LABEL_GAP = 26;
		var ROW_GAP = 8;

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
			.attr("role", "img")
			.attr("aria-label", "tree of " + agent.name + "'s skills and their dependencies");

		// the links group is created first so it paints behind the nodes; its
		// paths are bound once the columns are known.
		var linkGroup = svg.append("g")
			.attr("class", "links")
			.attr("fill", "none")
			.attr("stroke", "currentColor")
			.attr("stroke-opacity", 0.35)
			.attr("stroke-width", 1.5);

		var node = svg.append("g")
			.attr("class", "nodes")
			.selectAll("g")
			.data(root.descendants())
			.join("g")
			.attr("class", "node");

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

		// pass one: draw every label in full and measure it. Nothing is truncated
		// and nothing is estimated; the rendered width and height are what the
		// columns and rows are allocated from.
		node.append("text")
			.attr("dy", "0.31em")
			.attr("x", 8)
			.attr("paint-order", "stroke")
			.attr("stroke", "#fff")
			.attr("stroke-width", 3)
			.attr("stroke-linejoin", "round")
			.attr("font-family", "sans-serif")
			.attr("font-size", function (d) {
				return d.depth === 0 ? ROOT_FONT_PX : SKILL_FONT_PX;
			})
			.text(function (d) {
				return d.data.name;
			});

		// the row height comes from the measured height of a skill label, so the
		// rows are spaced by the font that is actually rendered.
		var skillFontHeight = SKILL_FONT_PX;
		node.each(function (d) {
			var text = this.querySelector("text");
			d.labelWidth = text.getComputedTextLength();
			if (d.depth > 0) skillFontHeight = text.getBBox().height;
		});

		var rowHeight = skillFontHeight + ROW_GAP;

		// pass two: a column is as wide as its widest label plus the gap, and the
		// columns are laid out end to end from those measured widths.
		var columnWidth = [];
		root.each(function (d) {
			var width = d.labelWidth + LABEL_GAP;
			columnWidth[d.depth] = Math.max(columnWidth[d.depth] || 0, width);
		});

		var columnX = [];
		var offset = ORIGIN_X;
		for (var depth = 0; depth < columnWidth.length; depth++) {
			columnX[depth] = offset;
			offset += columnWidth[depth];
		}

		var posX = function (d) {
			return columnX[d.depth];
		};

		var posY = function (d) {
			return ORIGIN_Y + d.row * rowHeight;
		};

		node.attr("transform", function (d) {
			return "translate(" + posX(d) + "," + posY(d) + ")";
		});

		// links: a step-before path, horizontal to the child's column then
		// vertical down to its row.
		var link = linkGroup.selectAll("path")
			.data(root.links())
			.join("path")
			.attr("d", function (d) {
				return "M" + posX(d.source) + "," + posY(d.source) + "H" + posX(d.target) + "V" + posY(d.target);
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

		// the svg is sized to the drawing at 1:1, so the labels render at their
		// stated pixel size and the stage scrolls whatever does not fit.
		var bounds = svg.node().getBBox();
		var width = Math.ceil(bounds.x + bounds.width + ORIGIN_X);
		var height = Math.ceil(bounds.y + bounds.height + ORIGIN_Y);
		svg.attr("width", width).attr("height", height).attr("viewBox", "0 0 " + width + " " + height);

		var caption = document.createElement("p");
		caption.className = "caption";
		caption.textContent = (root.descendants().length - 1) + " skills · " + root.leaves().length + " leaves";
		stage.appendChild(caption);
	};

	// wireSearch filters the list as the needle is typed. The filter is local,
	// so it never touches the url: only the selected agent is linkable.
	var wireSearch = function () {
		var input = document.querySelector("#agents-search");
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
