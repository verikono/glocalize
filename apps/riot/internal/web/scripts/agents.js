// agents browser client: lists every agent alphabetically, filters them from the
// search needle, and opens one into the main view via the ?agent= search param.
(function () {
	"use strict";

	// the narrowest the list column may be dragged, and the widest it may grow
	// as a fraction of the viewport.
	var MIN_WIDTH = 280;
	var MAX_FRACTION = 0.6;

	// the add-skill control: the radius of its round buttons, how far right of
	// the agent's name the button sits and how far above it, the padding inside
	// a popout, the height of a popout row, and the gap between the elements.
	var ADD_R = 10;
	var ADD_OFFSET_X = 25;
	var ADD_OFFSET_Y = 10;
	var POP_PAD = 8;
	var POP_ITEM_H = 12;
	var POP_GAP = 16;

	var PRIORITIES = ["LOW", "NORMAL", "HIGH", "ESSENTIAL"];

	// the line menu: its width and height, the gap from the line's elbow, and
	// the side of the checkbox inside it.
	var EDGE_POP_W = 150;
	var EDGE_POP_H = 26;
	var EDGE_POP_GAP = 10;
	var EDGE_BOX = 11;

	// the gap between a skill label's right edge and its add button.
	var SKILL_ADD_GAP = 14;

	var state = {
		agents: [],
		skills: [],
		filter: "",
		selected: "",
	};

	// view is the zoom and pan: a user-space point maps to screen as
	// k * point + (x, y). It outlives a redraw so the tree keeps its place.
	var view = { k: 1, x: 0, y: 0 };

	// edgeState is the line menu: whether it is open, which agent-skill line it
	// belongs to, the point on that line it is anchored to in tree space, and
	// whether its "use sub-tree" box is ticked. Tree space keeps the menu on the
	// same spot of the line as the tree is zoomed and panned.
	var edgeState = {
		open: false,
		sourceId: "",
		targetId: "",
		ux: 0,
		uy: 0,
		useSubtree: false,
	};

	// addState is the add-skill interface: whether it is open, which skills are
	// ticked, and the priority they will be bound at.
	var addState = {
		open: false,
		selected: new Set(),
		priority: "NORMAL",
	};

	// skillAddState is the add-dependency interface: whether it is open, which
	// skill it belongs to, which skills are ticked, and the search needle.
	var skillAddState = {
		open: false,
		skillId: "",
		selected: new Set(),
		needle: "",
	};

	// hoveredSkillId is the skill whose add button is showing, and hoverTimer
	// delays clearing it so the pointer can travel from the label to the button.
	var hoveredSkillId = "";
	var hoverTimer = null;
	var pointer = { x: -1, y: -1 };

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
			"{ Agents { id name description skills { priority skill { ...SkillTree } } } } " + dependencyFragment(DEPENDENCY_DEPTH)
		);
		if (result.error || !result.data) {
			console.error("Agents: failed getting data", result.error);
			return null;
		}
		return result.data.Agents || [];
	};

	// loadSkills fetches the whole skill catalogue, which the add list offers
	// from. It is separate from the agent query because the list needs every
	// skill, not just the ones an agent already has.
	var loadSkills = async () => {
		var result = await query("{ Skills { id name } }");
		if (result.error || !result.data) {
			console.error("Agents: failed getting skills", result.error);
			return null;
		}
		return result.data.Skills || [];
	};

	// byName orders the list alphabetically, case-insensitively, so the column
	// reads the same however the catalogue is stored.
	var byName = function (a, b) {
		return String(a.name || "").localeCompare(String(b.name || ""), undefined, { sensitivity: "base" });
	};

	// matchesTerm tests a haystack against the filter. A leading "!" reverses
	// the test, so "!kaltra" keeps everything without kaltra.
	var matchesTerm = function (haystack, term) {
		var raw = String(term || "").trim().toLowerCase();
		var negated = raw.charAt(0) === "!";
		var needle = negated ? raw.slice(1).trim() : raw;
		if (!needle) return true;
		var found = String(haystack || "").toLowerCase().indexOf(needle) !== -1;
		return negated ? !found : found;
	};

	var renderList = function () {
		var list = document.querySelector("#agents-items");
		if (!list) return;

		var term = state.filter.trim().toLowerCase();
		var agents = state.agents
			.filter(function (agent) {
				return matchesTerm(agent.name, term);
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

		// a different agent is a different tree, so the view and the line menu
		// start fresh.
		view = { k: 1, x: 0, y: 0 };
		edgeState.open = false;
		edgeState.sourceId = "";
		edgeState.targetId = "";
		edgeState.ux = 0;
		edgeState.uy = 0;
		edgeState.useSubtree = false;
		skillAddState.open = false;
		skillAddState.skillId = "";
		skillAddState.selected = new Set();
		skillAddState.needle = "";
		hoveredSkillId = "";

		var heading = document.querySelector("#agent-name");
		if (heading) heading.textContent = agent.name;

		renderList();
		renderAgentTree(agent);

		var url = new URL(window.location.href);
		url.searchParams.set("agent", agent.id);
		if (push) window.history.pushState({ agent: agent.id }, "", url.toString());
		else window.history.replaceState({ agent: agent.id }, "", url.toString());

		syncNavLinks();
	};

	// syncNavLinks carries the current context across to the skills UI, so
	// switching sections keeps it and the return trip restores it.
	var syncNavLinks = function () {
		var params = new URLSearchParams(window.location.search);
		var next = new URLSearchParams();
		["agent", "subject", "class"].forEach(function (key) {
			var value = String(params.get(key) || "").trim();
			if (value) next.set(key, value);
		});

		var query = next.toString();
		var suffix = query ? "?" + query : "";

		var skills = document.querySelector('.nav-link[href^="/skills/"]');
		if (skills) skills.href = "/skills/" + suffix;

		// the section's own link keeps its params too, so re-clicking it is a no-op.
		var agents = document.querySelector('.nav-link[href^="/agents/"]');
		if (agents) agents.href = "/agents/" + suffix;
	};

	// buildAgentTree turns an agent into a hierarchy rooted at the agent itself:
	// its skills are the first branch, and each skill's dependencies hang beneath
	// it. A skill is placed at most once, so a dependency shared by several
	// skills appears under the first one that claims it and a cycle cannot
	// recurse forever. Only the first branch carries a level: a dependency is
	// not bound to the agent, so it has no priority of its own.
	var buildAgentTree = function (agent) {
		var bindings = agent.skills || [];

		// index every skill we can see, including ones only reachable as a
		// dependency, so a nested stub can be replaced by its full record.
		var byId = new Map();
		var index = function (skill) {
			if (!skill || !skill.id) return;
			if (!byId.has(skill.id)) byId.set(skill.id, skill);
			(skill.dependencies || []).forEach(index);
		};
		bindings.forEach(function (binding) {
			index(binding.skill);
		});

		var placed = new Set();
		var build = function (skill, level) {
			placed.add(skill.id);
			var node = {
				id: skill.id,
				name: skill.name,
				description: skill.description || "",
				level: level || "",
				children: [],
			};
			(skill.dependencies || []).forEach(function (dep) {
				if (placed.has(dep.id)) return;
				node.children.push(build(byId.get(dep.id) || dep, ""));
			});
			node.children.sort(byName);
			return node;
		};

		var root = { id: agent.id, name: agent.name, description: agent.description || "", level: "", children: [] };
		bindings
			.slice()
			.sort(function (a, b) {
				return byName(a.skill, b.skill);
			})
			.forEach(function (binding) {
				if (!binding.skill || placed.has(binding.skill.id)) return;
				root.children.push(build(binding.skill, binding.priority));
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
			.attr("width", "100%")
			.attr("height", "100%")
			.attr("role", "img")
			.attr("aria-label", "tree of " + agent.name + "'s skills and their dependencies");

		// the tree lives in a group the zoom and pan transform, so the drawing
		// scales while the controls above it keep a fixed pixel size.
		var viewport = svg.append("g").attr("class", "viewport");

		// the links group is created first so it paints behind the nodes; its
		// paths are bound once the columns are known.
		var linkGroup = viewport.append("g")
			.attr("class", "links")
			.attr("fill", "none");

		var node = viewport.append("g")
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
				return d.data.level ? d.data.name + " (" + d.data.level + ")" : d.data.name;
			});

		// the root is the agent itself, so it is drawn as the tree's heading. The
		// class is applied before measuring, because it makes the label bold and
		// a bold label is wider than the width measured without it.
		node.filter(function (d) {
			return d.depth === 0;
		}).classed("agent-root", true);

		// sfh is a single skill's text height, measured from the label that was
		// actually drawn. The row pitch is that height plus a gap.
		var sfh = SKILL_FONT_PX;
		node.each(function (d) {
			var text = this.querySelector("text");
			d.labelWidth = text.getComputedTextLength();
			if (d.depth > 0) sfh = text.getBBox().height;
		});

		var rowHeight = sfh + ROW_GAP;

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
		// vertical down to its row. It starts at the parent label's right edge,
		// not the node, so the run is never hidden under the label's halo. Each
		// line is a group holding the drawn edge and a wider transparent path
		// that makes the thin line easy to hit.
		var linkPath = function (d) {
			var startX = posX(d.source) + 8 + d.source.labelWidth;
			return "M" + startX + "," + posY(d.source) + "H" + posX(d.target) + "V" + posY(d.target);
		};

		var link = linkGroup.selectAll("g.link")
			.data(root.links())
			.join("g")
			.attr("class", "link");

		var edge = link.append("path")
			.attr("class", "edge")
			.attr("stroke", "currentColor")
			.attr("stroke-opacity", 0.35)
			.attr("stroke-width", 1.5)
			.attr("d", linkPath);

		link.append("path")
			.attr("class", "hit")
			.attr("stroke", "transparent")
			.attr("stroke-width", 12)
			.attr("pointer-events", "stroke")
			.attr("d", linkPath);

		// hoveredPath is the ancestry of the node under the pointer, or null when
		// no node is hovered. It is kept so the line menu's own line can stay
		// bold while the pointer is elsewhere.
		var hoveredPath = null;

		// isActiveEdge reports whether a link is the one the line menu belongs to.
		var isActiveEdge = function (l) {
			return edgeState.open &&
				l.source.data.id === edgeState.sourceId &&
				l.target.data.id === edgeState.targetId;
		};

		// paintEdges draws the hover highlight and keeps the active line bold, so
		// the relationship the menu belongs to never fades while the menu is open.
		var paintEdges = function () {
			edge.attr("stroke-opacity", function (l) {
				if (isActiveEdge(l)) return 1;
				if (!hoveredPath) return 0.35;
				return hoveredPath.has(l.target) ? 0.95 : 0.12;
			}).attr("stroke-width", function (l) {
				if (isActiveEdge(l)) return 3;
				if (!hoveredPath) return 1.5;
				return hoveredPath.has(l.target) ? 2.5 : 1.5;
			});
		};

		// a left click on a line opens its menu at the pointer. The pointer is
		// recorded in tree space, so the menu stays on the same spot of the line
		// however the tree is later zoomed or panned.
		link.style("cursor", "pointer").on("click", function (event, d) {
			event.stopPropagation();
			var rect = svg.node().getBoundingClientRect();
			edgeState.open = true;
			edgeState.sourceId = d.source.data.id;
			edgeState.targetId = d.target.data.id;
			edgeState.ux = (event.clientX - rect.left - view.x) / view.k;
			edgeState.uy = (event.clientY - rect.top - view.y) / view.k;
			edgeState.useSubtree = false;
			paintEdges();
			renderOverlay();
		});

		// a right click on the active line deactivates the relationship: the menu
		// closes and the line returns to its resting weight.
		link.on("contextmenu", function (event, d) {
			event.preventDefault();
			event.stopPropagation();
			if (!isActiveEdge(d)) return;
			closeEdgeCallout();
		});

		node.append("title")
			.text(function (d) {
				return d.data.description ? d.data.name + "\n\n" + d.data.description : d.data.name;
			});

		// hovering a node highlights the path back to the agent; clicking a skill
		// opens it on the skills page. A skill also reveals its add button, which
		// is cleared on a delay so the pointer can travel from the label to it.
		node.on("mouseenter", function (event, d) {
			hoveredPath = new Set(d.ancestors());
			paintEdges();
			node.classed("active", function (n) {
				return hoveredPath.has(n);
			});
			if (d.depth === 0) return;
			clearTimeout(hoverTimer);
			if (hoveredSkillId !== d.data.id) {
				hoveredSkillId = d.data.id;
				renderOverlay();
			}
		}).on("mouseleave", function () {
			hoveredPath = null;
			paintEdges();
			node.classed("active", false);
			scheduleHoverClear();
		});

		node.style("cursor", function (d) {
			return d.depth === 0 ? "default" : "pointer";
		}).on("click", function (event, d) {
			if (d.depth === 0) return;
			event.stopPropagation();
			window.location.href = "/skills/?skill=" + encodeURIComponent(d.data.id);
		});

		// the overlay holds the add control and the line menu. It is never scaled,
		// so they keep their pixel size however far the tree is zoomed.
		var overlay = svg.append("g").attr("class", "overlay");

		// applyView writes the transform and re-anchors the overlay, which is
		// positioned in screen space from the tree.
		var applyView = function () {
			viewport.attr("transform", "translate(" + view.x + "," + view.y + ") scale(" + view.k + ")");
			renderOverlay();
		};

		// renderOverlay redraws everything that sits above the tree.
		var renderOverlay = function () {
			overlay.selectAll("*").remove();
			renderAddControl();
			renderSkillAddControl();
			renderEdgeCallout();
		};

		// scheduleHoverClear hides the skill's add button after a short delay, so
		// the pointer can cross the gap between the label and the button.
		var scheduleHoverClear = function () {
			clearTimeout(hoverTimer);
			hoverTimer = setTimeout(function () {
				if (skillAddState.open) return;
				if (hoveredSkillId === "") return;
				hoveredSkillId = "";
				renderOverlay();
			}, 220);
		};

		// renderSkillAddControl draws the round + button beside the hovered skill
		// and, when open, its dependency list with a search needle and the
		// confirm/cancel buttons. It offers only skills the tree does not show.
		var renderSkillAddControl = function () {
			if (hoveredSkillId === "" && !skillAddState.open) return;

			var skillId = skillAddState.open ? skillAddState.skillId : hoveredSkillId;
			var target = root.descendants().find(function (d) {
				return d.data.id === skillId;
			});
			if (!target) return;

			// the label's right edge in screen space, so the button tracks the tree
			// as it is zoomed and panned.
			var labelRight = view.x + view.k * (posX(target) + 8 + target.labelWidth);
			var cy = view.y + view.k * posY(target);
			var cx = labelRight + SKILL_ADD_GAP + ADD_R;
			var top = cy - ADD_R;

			var button = overlay.append("g").attr("class", "skill-add-button");
			button.append("circle")
				.attr("cx", cx)
				.attr("cy", cy)
				.attr("r", ADD_R)
				.attr("fill", "rgb(37, 37, 37)");
			button.append("text")
				.attr("x", cx)
				.attr("y", cy)
				.attr("dy", "0.35em")
				.attr("text-anchor", "middle")
				.attr("fill", "rgb(255, 255, 255)")
				.attr("font-family", "sans-serif")
				.attr("font-size", ADD_R * 0.96)
				.text("+");

			// entering the button keeps the skill hovered; leaving it schedules the
			// clear, so the button disappears once the pointer moves away.
			button.on("mouseenter", function () {
				clearTimeout(hoverTimer);
			}).on("mouseleave", function () {
				scheduleHoverClear();
			});

			button.style("cursor", "pointer").on("click", function (event) {
				event.stopPropagation();
				clearTimeout(hoverTimer);
				if (skillAddState.open) {
					closeSkillAddControl();
					return;
				}
				skillAddState.open = true;
				skillAddState.skillId = skillId;
				skillAddState.selected = new Set();
				skillAddState.needle = "";
				renderOverlay();
				svg.node().focus();
			});

			button.on("contextmenu", function (event) {
				event.preventDefault();
				event.stopPropagation();
				closeSkillAddControl();
			});

			if (!skillAddState.open) return;

			var visible = visibleSkillIDs(root);
			var term = skillAddState.needle.trim().toLowerCase();
			var offered = state.skills
				.filter(function (skill) {
					if (visible.has(skill.id) || skill.id === skillId) return false;
					if (!term) return true;
					return String(skill.name || "").toLowerCase().includes(term);
				})
				.sort(byName);

			// the list sits right of the button, with the needle as its first row.
			var listX = cx + ADD_R + POP_GAP;
			var listWidth = 260;
			var needleH = POP_ITEM_H + POP_PAD;
			var rowsTop = top + needleH;
			var listHeight = needleH + POP_PAD + Math.max(offered.length, 1) * POP_ITEM_H;

			var list = overlay.append("g").attr("class", "skill-add-list");
			list.append("rect")
				.attr("x", listX)
				.attr("y", top)
				.attr("width", listWidth)
				.attr("height", listHeight)
				.attr("rx", 4)
				.attr("fill", "rgb(255, 255, 255)")
				.attr("stroke", "rgb(37, 37, 37)")
				.attr("stroke-width", 1);

			// the needle row: a rule under it separates it from the list.
			list.append("rect")
				.attr("x", listX)
				.attr("y", top + needleH - 1)
				.attr("width", listWidth)
				.attr("height", 1)
				.attr("fill", "rgb(190, 190, 190)");

			list.append("text")
				.attr("x", listX + POP_PAD)
				.attr("y", top + POP_PAD)
				.attr("dy", "0.8em")
				.attr("font-family", "sans-serif")
				.attr("font-size", POP_ITEM_H)
				.attr("fill", skillAddState.needle ? "rgb(37, 37, 37)" : "rgb(150, 150, 150)")
				.text(skillAddState.needle || "search skills");

			if (offered.length === 0) {
				list.append("text")
					.attr("x", listX + POP_PAD)
					.attr("y", rowsTop + POP_PAD)
					.attr("dy", "0.8em")
					.attr("font-family", "sans-serif")
					.attr("font-size", POP_ITEM_H)
					.attr("fill", "rgb(110, 110, 110)")
					.text(term ? "no matching skills" : "every skill is already shown");
			}

			offered.forEach(function (skill, index) {
				var rowY = rowsTop + POP_PAD + index * POP_ITEM_H;
				var row = list.append("g").attr("class", "skill-add-item");
				var selected = skillAddState.selected.has(skill.id);

				row.append("rect")
					.attr("x", listX)
					.attr("y", rowY)
					.attr("width", listWidth)
					.attr("height", POP_ITEM_H)
					.attr("fill", selected ? "rgb(37, 37, 37)" : "transparent");

				row.append("text")
					.attr("x", listX + POP_PAD)
					.attr("y", rowY)
					.attr("dy", "0.8em")
					.attr("font-family", "sans-serif")
					.attr("font-size", POP_ITEM_H)
					.attr("fill", selected ? "rgb(255, 255, 255)" : "rgb(37, 37, 37)")
					.text(skill.name);

				row.style("cursor", "pointer").on("click", function (event) {
					event.stopPropagation();
					if (skillAddState.selected.has(skill.id)) skillAddState.selected.delete(skill.id);
					else skillAddState.selected.add(skill.id);
					renderOverlay();
				});
			});

			// the confirm and cancel buttons, right of the list. There is no
			// priority list here: a dependency has no priority of its own.
			var confirmX = listX + listWidth + POP_GAP + ADD_R;
			var cancelX = confirmX + ADD_R * 2 + POP_GAP + ADD_R;

			var confirm = overlay.append("g").attr("class", "skill-add-confirm");
			confirm.append("circle")
				.attr("cx", confirmX)
				.attr("cy", cy)
				.attr("r", ADD_R)
				.attr("fill", "rgb(40, 150, 70)");
			confirm.append("path")
				.attr("d", "M" + (confirmX - ADD_R * 0.44) + "," + cy +
					"L" + (confirmX - ADD_R * 0.12) + "," + (cy + ADD_R * 0.32) +
					"L" + (confirmX + ADD_R * 0.48) + "," + (cy - ADD_R * 0.32))
				.attr("fill", "none")
				.attr("stroke", "rgb(255, 255, 255)")
				.attr("stroke-width", ADD_R * 0.16)
				.attr("stroke-linecap", "round")
				.attr("stroke-linejoin", "round");
			confirm.style("cursor", "pointer").on("click", function (event) {
				event.stopPropagation();
				addSelectedDependencies(skillId);
			});

			var cancel = overlay.append("g").attr("class", "skill-add-cancel");
			cancel.append("circle")
				.attr("cx", cancelX)
				.attr("cy", cy)
				.attr("r", ADD_R)
				.attr("fill", "rgb(190, 50, 50)");
			cancel.append("path")
				.attr("d", "M" + (cancelX - ADD_R * 0.36) + "," + (cy - ADD_R * 0.36) +
					"L" + (cancelX + ADD_R * 0.36) + "," + (cy + ADD_R * 0.36) +
					"M" + (cancelX + ADD_R * 0.36) + "," + (cy - ADD_R * 0.36) +
					"L" + (cancelX - ADD_R * 0.36) + "," + (cy + ADD_R * 0.36))
				.attr("fill", "none")
				.attr("stroke", "rgb(255, 255, 255)")
				.attr("stroke-width", ADD_R * 0.16)
				.attr("stroke-linecap", "round");
			cancel.style("cursor", "pointer").on("click", function (event) {
				event.stopPropagation();
				closeSkillAddControl();
			});
		};

		// renderAddControl draws the round + button beside the agent's name and,
		// when open, the skill list, the priority list and the confirm/cancel
		// buttons. Every element's top edge is aligned to the button's top edge.
		var renderAddControl = function () {
			// the agent's name in screen space, so the control tracks the tree as
			// it is zoomed and panned.
			var nameX = view.x + view.k * posX(root);
			var nameY = view.y + view.k * posY(root);
			var nameRight = nameX + view.k * (8 + root.labelWidth);

			// the button sits ADD_OFFSET_X right of the name and ADD_OFFSET_Y
			// above it.
			var cx = nameRight + ADD_OFFSET_X + ADD_R;
			var cy = nameY - ADD_OFFSET_Y;
			var top = cy - ADD_R;

			var button = overlay.append("g").attr("class", "add-button");
			button.append("circle")
				.attr("cx", cx)
				.attr("cy", cy)
				.attr("r", ADD_R)
				.attr("fill", "rgb(37, 37, 37)");
			button.append("text")
				.attr("x", cx)
				.attr("y", cy)
				.attr("dy", "0.35em")
				.attr("text-anchor", "middle")
				.attr("fill", "rgb(255, 255, 255)")
				.attr("font-family", "sans-serif")
				.attr("font-size", ADD_R * 0.96)
				.text("+");

			button.style("cursor", "pointer").on("click", function (event) {
				event.stopPropagation();
				if (addState.open) {
					closeAddControl();
					return;
				}
				addState.open = true;
				addState.selected = new Set();
				addState.priority = "NORMAL";
				renderAgentTree(agent);
			});

			// a right click on the button cancels the interaction.
			button.on("contextmenu", function (event) {
				event.preventDefault();
				event.stopPropagation();
				closeAddControl();
			});

			if (!addState.open) return;

			var visible = visibleSkillIDs(root);
			var offered = state.skills
				.filter(function (skill) {
					return !visible.has(skill.id);
				})
				.sort(byName);

			// the skill list popout, immediately right of the button.
			var listX = cx + ADD_R + POP_GAP;
			var listWidth = 260;
			var listHeight = POP_PAD * 2 + Math.max(offered.length, 1) * POP_ITEM_H;

			var list = overlay.append("g").attr("class", "add-list");
			list.append("rect")
				.attr("x", listX)
				.attr("y", top)
				.attr("width", listWidth)
				.attr("height", listHeight)
				.attr("rx", 4)
				.attr("fill", "rgb(255, 255, 255)")
				.attr("stroke", "rgb(37, 37, 37)")
				.attr("stroke-width", 1);

			if (offered.length === 0) {
				list.append("text")
					.attr("x", listX + POP_PAD)
					.attr("y", top + POP_PAD)
					.attr("dy", "0.8em")
					.attr("font-family", "sans-serif")
					.attr("font-size", POP_ITEM_H)
					.attr("fill", "rgb(110, 110, 110)")
					.text("every skill is already shown");
			}

			offered.forEach(function (skill, index) {
				var rowY = top + POP_PAD + index * POP_ITEM_H;
				var row = list.append("g").attr("class", "add-item");
				var selected = addState.selected.has(skill.id);

				row.append("rect")
					.attr("x", listX)
					.attr("y", rowY)
					.attr("width", listWidth)
					.attr("height", POP_ITEM_H)
					.attr("fill", selected ? "rgb(37, 37, 37)" : "transparent");

				row.append("text")
					.attr("x", listX + POP_PAD)
					.attr("y", rowY)
					.attr("dy", "0.8em")
					.attr("font-family", "sans-serif")
					.attr("font-size", POP_ITEM_H)
					.attr("fill", selected ? "rgb(255, 255, 255)" : "rgb(37, 37, 37)")
					.text(skill.name);

				row.style("cursor", "pointer").on("click", function (event) {
					event.stopPropagation();
					if (addState.selected.has(skill.id)) addState.selected.delete(skill.id);
					else addState.selected.add(skill.id);
					renderAgentTree(agent);
				});
			});

			// the priority list and the confirm/cancel buttons only appear once at
			// least one skill is ticked.
			if (addState.selected.size === 0) return;

			var priorityX = listX + listWidth + POP_GAP;
			var priorityWidth = 120;
			var priorityHeight = POP_PAD * 2 + PRIORITIES.length * POP_ITEM_H;

			var priorities = overlay.append("g").attr("class", "add-priorities");
			priorities.append("rect")
				.attr("x", priorityX)
				.attr("y", top)
				.attr("width", priorityWidth)
				.attr("height", priorityHeight)
				.attr("rx", 4)
				.attr("fill", "rgb(255, 255, 255)")
				.attr("stroke", "rgb(37, 37, 37)")
				.attr("stroke-width", 1);

			PRIORITIES.forEach(function (priority, index) {
				var rowY = top + POP_PAD + index * POP_ITEM_H;
				var selected = addState.priority === priority;
				var row = priorities.append("g").attr("class", "add-priority");

				row.append("rect")
					.attr("x", priorityX)
					.attr("y", rowY)
					.attr("width", priorityWidth)
					.attr("height", POP_ITEM_H)
					.attr("fill", selected ? "rgb(37, 37, 37)" : "transparent");

				row.append("text")
					.attr("x", priorityX + POP_PAD)
					.attr("y", rowY)
					.attr("dy", "0.8em")
					.attr("font-family", "sans-serif")
					.attr("font-size", POP_ITEM_H)
					.attr("fill", selected ? "rgb(255, 255, 255)" : "rgb(37, 37, 37)")
					.text(priority);

				row.style("cursor", "pointer").on("click", function (event) {
					event.stopPropagation();
					addState.priority = priority;
					renderAgentTree(agent);
				});
			});

			// the confirm and cancel buttons, right of the priority list. Their
			// glyphs are drawn as fractions of the radius, so they scale with it.
			var confirmX = priorityX + priorityWidth + POP_GAP + ADD_R;
			var cancelX = confirmX + ADD_R * 2 + POP_GAP + ADD_R;

			var confirm = overlay.append("g").attr("class", "add-confirm");
			confirm.append("circle")
				.attr("cx", confirmX)
				.attr("cy", cy)
				.attr("r", ADD_R)
				.attr("fill", "rgb(40, 150, 70)");
			confirm.append("path")
				.attr("d", "M" + (confirmX - ADD_R * 0.44) + "," + cy +
					"L" + (confirmX - ADD_R * 0.12) + "," + (cy + ADD_R * 0.32) +
					"L" + (confirmX + ADD_R * 0.48) + "," + (cy - ADD_R * 0.32))
				.attr("fill", "none")
				.attr("stroke", "rgb(255, 255, 255)")
				.attr("stroke-width", ADD_R * 0.16)
				.attr("stroke-linecap", "round")
				.attr("stroke-linejoin", "round");
			confirm.style("cursor", "pointer").on("click", function (event) {
				event.stopPropagation();
				bindSelectedSkills(agent);
			});

			var cancel = overlay.append("g").attr("class", "add-cancel");
			cancel.append("circle")
				.attr("cx", cancelX)
				.attr("cy", cy)
				.attr("r", ADD_R)
				.attr("fill", "rgb(190, 50, 50)");
			cancel.append("path")
				.attr("d", "M" + (cancelX - ADD_R * 0.36) + "," + (cy - ADD_R * 0.36) +
					"L" + (cancelX + ADD_R * 0.36) + "," + (cy + ADD_R * 0.36) +
					"M" + (cancelX + ADD_R * 0.36) + "," + (cy - ADD_R * 0.36) +
					"L" + (cancelX - ADD_R * 0.36) + "," + (cy + ADD_R * 0.36))
				.attr("fill", "none")
				.attr("stroke", "rgb(255, 255, 255)")
				.attr("stroke-width", ADD_R * 0.16)
				.attr("stroke-linecap", "round");
			cancel.style("cursor", "pointer").on("click", function (event) {
				event.stopPropagation();
				closeAddControl();
			});

			// a right click on any part of the open interface cancels it.
			overlay.on("contextmenu", function (event) {
				event.preventDefault();
				event.stopPropagation();
				closeAddControl();
			});
		};

		// closeEdgeCallout deactivates the relationship: the menu closes and the
		// line it belonged to returns to its resting weight.
		var closeEdgeCallout = function () {
			if (!edgeState.open) return;
			edgeState.open = false;
			edgeState.sourceId = "";
			edgeState.targetId = "";
			edgeState.useSubtree = false;
			paintEdges();
			renderOverlay();
		};

		// closeSkillAddControl drops the dependency interface and its button.
		var closeSkillAddControl = function () {
			skillAddState.open = false;
			skillAddState.skillId = "";
			skillAddState.selected = new Set();
			skillAddState.needle = "";
			hoveredSkillId = "";
			clearTimeout(hoverTimer);
			renderOverlay();
		};

		// addSelectedDependencies links the ticked skills to the hovered skill as
		// its dependencies, then reloads the catalogue so the tree shows them.
		var addSelectedDependencies = async function (skillId) {
			var dependencyIds = Array.from(skillAddState.selected);
			if (dependencyIds.length === 0) return;

			var result = await query(
				"mutation AddDeps($skillId: String!, $dependencyIds: [String!]!) { AddSkillDependencies(skillId: $skillId, dependencyIds: $dependencyIds) { id name } }",
				{ skillId: skillId, dependencyIds: dependencyIds }
			);
			if (result.error) {
				console.error("Agents: failed adding dependencies", result.error);
				return;
			}

			closeSkillAddControl();

			var agents = await loadAgents();
			if (agents) state.agents = agents;
			var refreshed = state.agents.find(function (candidate) {
				return candidate.id === state.selected;
			});
			if (refreshed) renderAgentTree(refreshed);
		};

		// renderEdgeCallout draws the menu for the clicked agent-skill line: one
		// "use sub-tree" row with a checkbox, hung just below and right of the
		// anchor point so it never covers the line it belongs to. The anchor is in
		// tree space, so the menu follows the line as the tree is zoomed and panned
		// while the menu itself keeps its pixel size.
		var renderEdgeCallout = function () {
			if (!edgeState.open) return;

			var x = view.x + view.k * edgeState.ux + EDGE_POP_GAP;
			var y = view.y + view.k * edgeState.uy + EDGE_POP_GAP;

			var callout = overlay.append("g").attr("class", "edge-callout");

			callout.append("rect")
				.attr("x", x)
				.attr("y", y)
				.attr("width", EDGE_POP_W)
				.attr("height", EDGE_POP_H)
				.attr("rx", 4)
				.attr("fill", "rgb(255, 255, 255)")
				.attr("stroke", "rgb(37, 37, 37)")
				.attr("stroke-width", 1);

			var boxX = x + POP_PAD;
			var boxY = y + (EDGE_POP_H - EDGE_BOX) / 2;

			callout.append("rect")
				.attr("x", boxX)
				.attr("y", boxY)
				.attr("width", EDGE_BOX)
				.attr("height", EDGE_BOX)
				.attr("rx", 2)
				.attr("fill", edgeState.useSubtree ? "rgb(37, 37, 37)" : "rgb(255, 255, 255)")
				.attr("stroke", "rgb(37, 37, 37)")
				.attr("stroke-width", 1);

			if (edgeState.useSubtree) {
				callout.append("path")
					.attr("d", "M" + (boxX + EDGE_BOX * 0.22) + "," + (boxY + EDGE_BOX * 0.52) +
						"L" + (boxX + EDGE_BOX * 0.42) + "," + (boxY + EDGE_BOX * 0.74) +
						"L" + (boxX + EDGE_BOX * 0.8) + "," + (boxY + EDGE_BOX * 0.26))
					.attr("fill", "none")
					.attr("stroke", "rgb(255, 255, 255)")
					.attr("stroke-width", 1.6)
					.attr("stroke-linecap", "round")
					.attr("stroke-linejoin", "round");
			}

			callout.append("text")
				.attr("x", boxX + EDGE_BOX + POP_PAD)
				.attr("y", y + EDGE_POP_H / 2)
				.attr("dy", "0.35em")
				.attr("font-family", "sans-serif")
				.attr("font-size", 12)
				.attr("fill", "rgb(37, 37, 37)")
				.text("use sub-tree");

			callout.style("cursor", "pointer").on("click", function (event) {
				event.stopPropagation();
				edgeState.useSubtree = !edgeState.useSubtree;
				renderOverlay();
			});

			// a right click on the callout deactivates the relationship too.
			callout.on("contextmenu", function (event) {
				event.preventDefault();
				event.stopPropagation();
				closeEdgeCallout();
			});
		};

		// wireZoomPan binds wheel zoom about the pointer and middle-drag pan. The
		// overlay is re-anchored on every change, so the control stays beside the
		// agent's name at a constant size.
		var wireZoomPan = function () {
			var svgNode = svg.node();

			svgNode.addEventListener("wheel", function (event) {
				event.preventDefault();

				var rect = svgNode.getBoundingClientRect();
				var px = event.clientX - rect.left;
				var py = event.clientY - rect.top;

				var k = Math.min(Math.max(view.k * Math.exp(-event.deltaY * 0.0015), 0.2), 8);
				if (k === view.k) return;

				// keep the user-space point under the pointer fixed.
				var userX = (px - view.x) / view.k;
				var userY = (py - view.y) / view.k;
				view.k = k;
				view.x = px - k * userX;
				view.y = py - k * userY;
				applyView();
			}, { passive: false });

			var panning = false;
			var panStart = null;

			var onPanMove = function (event) {
				if (!panning || !panStart) return;
				event.preventDefault();
				view.x = panStart.x + (event.clientX - panStart.cx);
				view.y = panStart.y + (event.clientY - panStart.cy);
				applyView();
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
				// button 1 is the middle mouse button.
				if (event.button !== 1) return;
				event.preventDefault();

				panning = true;
				panStart = { cx: event.clientX, cy: event.clientY, x: view.x, y: view.y };
				svgNode.classList.add("panning");
				svgNode.setPointerCapture?.(event.pointerId);
				window.addEventListener("pointermove", onPanMove);
				window.addEventListener("pointerup", onPanUp);
				window.addEventListener("pointercancel", onPanUp);
			});

			// suppress the browser's middle-click autoscroll and paste behaviours.
			svgNode.addEventListener("auxclick", function (event) {
				if (event.button === 1) event.preventDefault();
			});
			svgNode.addEventListener("mousedown", function (event) {
				if (event.button === 1) event.preventDefault();
			});

			// a click on empty space dismisses the line menu and the dependency
			// interface.
			svgNode.addEventListener("click", function () {
				closeEdgeCallout();
				if (skillAddState.open) closeSkillAddControl();
			});

			// the dependency list's needle is typed straight into the svg, which
			// takes focus while the list is open.
			svgNode.setAttribute("tabindex", "0");
			svgNode.addEventListener("keydown", function (event) {
				if (!skillAddState.open) return;
				if (event.key === "Escape") {
					event.preventDefault();
					closeSkillAddControl();
					return;
				}
				if (event.key === "Backspace") {
					event.preventDefault();
					skillAddState.needle = skillAddState.needle.slice(0, -1);
					renderOverlay();
					return;
				}
				if (event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey) {
					event.preventDefault();
					skillAddState.needle += event.key;
					renderOverlay();
				}
			});
		};

		applyView();
		wireZoomPan();

		var caption = document.createElement("p");
		caption.className = "caption";
		caption.textContent = (root.descendants().length - 1) + " skills · " + root.leaves().length + " leaves";
		stage.appendChild(caption);
	};

	// visibleSkillIDs is the set of skills already drawn in the tree. The add
	// list filters them out, so it only ever offers skills the agent lacks.
	var visibleSkillIDs = function (root) {
		var ids = new Set();
		root.each(function (d) {
			if (d.depth > 0) ids.add(d.data.id);
		});
		return ids;
	};

	// closeAddControl drops the interface and redraws the tree without it.
	var closeAddControl = function () {
		addState.open = false;
		addState.selected = new Set();
		addState.priority = "NORMAL";
		var agent = state.agents.find(function (candidate) {
			return candidate.id === state.selected;
		});
		if (agent) renderAgentTree(agent);
	};

	// bindSelectedSkills attaches the ticked skills to the agent at the chosen
	// priority, then reloads the catalogue so the tree shows them.
	var bindSelectedSkills = async function (agent) {
		var skillIds = Array.from(addState.selected);
		if (skillIds.length === 0) return;

		var result = await query(
			"mutation Bind($agentId: String!, $skillIds: [String!]!, $priority: String) { BindAgentSkills(agentId: $agentId, skillIds: $skillIds, priority: $priority) { id name } }",
			{ agentId: agent.id, skillIds: skillIds, priority: addState.priority }
		);
		if (result.error) {
			console.error("Agents: failed binding skills", result.error);
			return;
		}

		addState.open = false;
		addState.selected = new Set();
		addState.priority = "NORMAL";

		var agents = await loadAgents();
		if (agents) state.agents = agents;
		var refreshed = state.agents.find(function (candidate) {
			return candidate.id === agent.id;
		});
		if (refreshed) renderAgentTree(refreshed);
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

		// the add list offers from the whole catalogue, so it is loaded alongside
		// the agents rather than on demand.
		var skills = await loadSkills();
		if (skills) state.skills = skills;

		// ?agent=<identifier> opens straight onto that agent.
		var params = new URLSearchParams(window.location.search);
		var wanted = String(params.get("agent") || "").trim();
		if (wanted) selectAgent(wanted, { push: false });
		else syncNavLinks();

		window.addEventListener("popstate", function () {
			var next = new URLSearchParams(window.location.search);
			var id = String(next.get("agent") || "").trim();
			if (id) selectAgent(id, { push: false });
			else syncNavLinks();
		});
	};

	document.addEventListener("DOMContentLoaded", function () {
		main();
	});
})();
