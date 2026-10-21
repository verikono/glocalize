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

	var loadAgents = async () => {
		var result = await query("{ Agents { id name description } }");
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

		var url = new URL(window.location.href);
		url.searchParams.set("agent", agent.id);
		if (push) window.history.pushState({ agent: agent.id }, "", url.toString());
		else window.history.replaceState({ agent: agent.id }, "", url.toString());
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
