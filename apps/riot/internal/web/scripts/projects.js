// projects browser client: lists every project alphabetically, filters them from the search needle, and opens one into the main view via the ?project= search param.
(function () {
	"use strict";

	// the narrowest the list column may be dragged, and the widest it may grow as a fraction of the viewport.
	var MIN_WIDTH = 280;
	var MAX_FRACTION = 0.6;

	// the indent applied per tree depth, in pixels.
	var INDENT_PX = 16;

	var state = {
		projects: [],
		filter: "",
		selected: "",
		collapsed: new Set(),
		loaded: null,
		isNew: false,
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

	// SKILL_DEPTH is how many levels of dependencies the project tree asks for: the project's skills, their dependencies, and those dependencies' dependencies.
	var SKILL_DEPTH = 2;

	// dependencyFragment returns a GraphQL fragment selecting a skill and its dependencies to the given depth. Each level nests the one before it, so a single query returns the whole chain rather than one level of it.
	var dependencyFragment = function (depth) {
		var fields = "id name description";
		for (var i = 0; i < depth; i++) {
			fields = "id name description dependencies { " + fields + " }";
		}
		return "fragment SkillTree on Skill { " + fields + " }";
	};

	var loadProjects = async () => {
		var result = await query("{ Projects { name location skillsPath agentsPath agents { id name description skills { priority skill { ...SkillTree } } } } } " + dependencyFragment(SKILL_DEPTH));
		if (result.error || !result.data) {
			console.error("Projects: failed getting data", result.error);
			return null;
		}
		return result.data.Projects || [];
	};

	// byName orders the list alphabetically, case-insensitively, so the column reads the same however the catalogue is stored.
	var byName = function (a, b) {
		return String(a.name || "").localeCompare(String(b.name || ""), undefined, { sensitivity: "base" });
	};

	// matchesTerm tests a haystack against the filter. A leading "!" reverses the test, so "!kaltra" keeps everything without kaltra.
	var matchesTerm = function (haystack, term) {
		var raw = String(term || "").trim().toLowerCase();
		var negated = raw.charAt(0) === "!";
		var needle = negated ? raw.slice(1).trim() : raw;
		if (!needle) return true;
		var found = String(haystack || "").toLowerCase().indexOf(needle) !== -1;
		return negated ? !found : found;
	};

	// findProject returns the project with the given name, case-insensitively, or null.
	var findProject = function (name) {
		var needle = String(name || "").trim().toLowerCase();
		if (!needle) return null;
		return (
			state.projects.find(function (project) {
				return String(project.name || "").toLowerCase() === needle;
			}) || null
		);
	};

	// fillProjectForm writes a project into the form; the name is read-only when updating because it is the primary key and is never renamed.
	var fillProjectForm = function (project) {
		var record = project || {};

		var title = document.querySelector("#project-form-title");
		if (title) title.textContent = state.isNew ? "New project" : "Update " + (record.name || "");

		var revert = document.querySelector("#project-form-revert");
		if (revert) revert.textContent = state.isNew ? "Cancel" : "Revert";

		var save = document.querySelector("#project-form-save");
		if (save) save.textContent = state.isNew ? "Create" : "Save";

		var name = document.querySelector("#project-form-name");
		if (name) {
			name.value = record.name || "";
			name.readOnly = !state.isNew;
		}

		var location = document.querySelector("#project-form-location");
		if (location) location.value = record.location || "";

		var skillsPath = document.querySelector("#project-form-skills-path");
		if (skillsPath) skillsPath.value = record.skillsPath || "";

		var agentsPath = document.querySelector("#project-form-agents-path");
		if (agentsPath) agentsPath.value = record.agentsPath || "";
	};

	// setProjectStatus writes a short message into the form's status line.
	var setProjectStatus = function (message, isError) {
		var status = document.querySelector("#project-form-status");
		if (!status) return;
		status.textContent = message;
		status.classList.toggle("errored", !!isError);
	};

	var startNewProject = function () {
		state.isNew = true;
		state.loaded = null;
		fillProjectForm(null);
		setProjectStatus("new project");
		var name = document.querySelector("#project-form-name");
		if (name) name.focus();
	};

	// revertProjectForm cancels a create back to the focused project, or restores the loaded record's values.
	var revertProjectForm = function () {
		if (state.isNew) {
			state.isNew = false;
			state.loaded = findProject(state.selected) || findProject("all");
			fillProjectForm(state.loaded);
			setProjectStatus("cancelled");
			return;
		}
		fillProjectForm(state.loaded);
		setProjectStatus("reverted");
	};

	// saveProject posts the form through the SaveProject mutation, then reloads the list and redraws the tree for the saved project.
	var saveProject = async function () {
		var nameEl = document.querySelector("#project-form-name");
		var locationEl = document.querySelector("#project-form-location");
		var skillsEl = document.querySelector("#project-form-skills-path");
		var agentsEl = document.querySelector("#project-form-agents-path");

		var name = nameEl ? nameEl.value.trim() : "";
		var location = locationEl ? locationEl.value.trim() : "";
		var skillsPath = skillsEl ? skillsEl.value.trim() : "";
		var agentsPath = agentsEl ? agentsEl.value.trim() : "";

		if (!name) {
			setProjectStatus("name is required", true);
			return;
		}
		if (state.isNew && findProject(name)) {
			setProjectStatus("a project named " + name + " already exists", true);
			return;
		}

		var input = { name: name, location: location, skillsPath: skillsPath, agentsPath: agentsPath, isNew: state.isNew };

		setProjectStatus("saving…");
		var result = await query("mutation SaveProject($input: ProjectInput!) { SaveProject(input: $input) { name location skillsPath agentsPath } }", { input: input });
		if (result.error || !result.data || !result.data.SaveProject) {
			var failure = result.error;
			var message = failure && failure.message ? failure.message : Array.isArray(failure) && failure[0] && failure[0].message ? failure[0].message : String(failure || "unknown error");
			setProjectStatus("save failed: " + message, true);
			return;
		}

		var saved = result.data.SaveProject;
		var wasNew = state.isNew;
		state.loaded = saved;
		state.isNew = false;
		setProjectStatus("saved");

		var projects = await loadProjects();
		if (projects) state.projects = projects;
		renderList();

		// a create focuses the new project; an update keeps the focused one.
		var focusName = wasNew ? saved.name : state.selected || saved.name;
		selectProject(focusName, { push: false });
	};

	var renderList = function () {
		var list = document.querySelector("#projects-items");
		if (!list) return;

		var term = state.filter.trim().toLowerCase();
		var projects = state.projects
			.filter(function (project) {
				return String(project.name || "").toLowerCase() !== "all";
			})
			.filter(function (project) {
				return matchesTerm(project.name, term);
			})
			.sort(byName);

		list.textContent = "";

		if (projects.length === 0) {
			var empty = document.createElement("div");
			empty.className = "projects-empty";
			empty.textContent = term ? "no matching projects" : "no projects";
			list.appendChild(empty);
			return;
		}

		projects.forEach(function (project) {
			var button = document.createElement("button");
			button.type = "button";
			button.className = "project-item";
			button.dataset.project = project.name;
			button.textContent = project.name;
			button.title = project.name;
			button.classList.toggle("active", project.name === state.selected);
			button.addEventListener("click", function () {
				selectProject(project.name);
			});
			list.appendChild(button);
		});
	};

	// selectProject opens a project in the main view and records it in the url, so the current project survives a reload and can be linked to. An empty or unknown name shows the "all" project.
	var selectProject = function (name, options) {
		var push = !options || options.push !== false;
		var project = findProject(name);
		if (project && String(project.name).toLowerCase() === "all") project = null;

		state.selected = project ? project.name : "";

		// the form always shows something editable: the focused project, or "all" when nothing is focused.
		state.isNew = false;
		state.loaded = project || findProject("all");
		fillProjectForm(state.loaded);

		var heading = document.querySelector("#project-name");
		if (heading) heading.textContent = project ? project.name : "All";

		renderList();
		renderProjectTree(project || findProject("all"));

		var url = new URL(window.location.href);
		if (project) url.searchParams.set("project", project.name);
		else url.searchParams.delete("project");
		if (push) window.history.pushState({ project: state.selected }, "", url.toString());
		else window.history.replaceState({ project: state.selected }, "", url.toString());
	};

	// buildProjectTree flattens a project into rows: each agent, its bound skills, and two levels of their dependencies. A skill is placed at most once within one agent's subtree, so a shared dependency appears under the first skill that claims it and a cycle cannot recurse forever.
	var buildProjectTree = function (project) {
		var rows = [];
		var agents = (project.agents || []).slice().sort(byName);

		agents.forEach(function (agent) {
			var collapsed = state.collapsed.has(agent.id);
			rows.push({ kind: "agent", depth: 0, id: agent.id, label: agent.name, collapsed: collapsed });
			if (collapsed) return;

			var bindings = (agent.skills || []).slice().sort(function (a, b) {
				return byName(a.skill, b.skill);
			});
			if (bindings.length === 0) {
				rows.push({ kind: "empty", depth: 1, label: "no skills" });
				return;
			}

			var placed = new Set();
			var walk = function (skill, depth, priority) {
				if (!skill || !skill.id || placed.has(skill.id)) return;
				placed.add(skill.id);
				rows.push({ kind: "skill", depth: depth, id: skill.id, label: priority ? skill.name + " (" + priority + ")" : skill.name });
				if (depth >= 3) return;
				(skill.dependencies || [])
					.slice()
					.sort(byName)
					.forEach(function (dep) {
						walk(dep, depth + 1, "");
					});
			};

			bindings.forEach(function (binding) {
				walk(binding.skill, 1, binding.priority);
			});
		});

		return rows;
	};

	// renderProjectTree draws the project as an indented list: the agent at the root, its skills beneath it, and their dependencies below those.
	var renderProjectTree = function (project) {
		var tree = document.querySelector("#project-tree");
		if (!tree) return;

		tree.textContent = "";

		if (!project) {
			var missing = document.createElement("p");
			missing.className = "projects-empty";
			missing.textContent = "no projects";
			tree.appendChild(missing);
			return;
		}

		var rows = buildProjectTree(project);
		if (rows.length === 0) {
			var empty = document.createElement("p");
			empty.className = "projects-empty";
			empty.textContent = "this project has no agents";
			tree.appendChild(empty);
			return;
		}

		rows.forEach(function (row) {
			var element = document.createElement("div");
			element.className = "tree-row";
			element.style.paddingLeft = row.depth * INDENT_PX + "px";

			var marker = document.createElement("span");
			marker.className = "tree-chevron";
			element.appendChild(marker);

			var label = document.createElement("span");
			label.textContent = row.label;
			element.appendChild(label);

			if (row.kind === "agent") {
				element.classList.add("tree-agent");
				marker.textContent = row.collapsed ? "▸" : "▾";
				element.addEventListener("click", function () {
					if (state.collapsed.has(row.id)) state.collapsed.delete(row.id);
					else state.collapsed.add(row.id);
					renderProjectTree(project);
				});
			} else if (row.kind === "skill") {
				element.classList.add("tree-skill");
				element.addEventListener("click", function () {
					window.location.href = "/skills/?skill=" + encodeURIComponent(row.id);
				});
			} else {
				element.classList.add("tree-empty");
			}

			tree.appendChild(element);
		});
	};

	var wireProjectForm = function () {
		var form = document.querySelector("#project-form");
		if (!form) return;

		var newButton = document.querySelector("#project-form-new");
		if (newButton) newButton.addEventListener("click", startNewProject);

		var revertButton = document.querySelector("#project-form-revert");
		if (revertButton) revertButton.addEventListener("click", revertProjectForm);

		form.addEventListener("submit", function (event) {
			event.preventDefault();
			saveProject();
		});
	};

	var wireSearch = function () {
		var input = document.querySelector("#projects-search");
		if (!input) return;
		input.addEventListener("input", function () {
			state.filter = input.value;
			renderList();
		});
	};

	// wireResizer makes the separator drag the boundary between the list and the main view. Only the list's width is set; the main view takes the rest.
	var wireResizer = function () {
		var handle = document.querySelector("#projects-resizer");
		var list = document.querySelector("#projects-list");
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
		wireProjectForm();

		var projects = await loadProjects();
		if (!projects) {
			var list = document.querySelector("#projects-items");
			if (list) {
				list.textContent = "";
				var p = document.createElement("p");
				p.className = "errored";
				p.textContent = "failed to load projects";
				list.appendChild(p);
			}
			return;
		}

		state.projects = projects;
		renderList();

		// ?project=<name> opens straight onto that project; anything else shows "all".
		var params = new URLSearchParams(window.location.search);
		selectProject(String(params.get("project") || "").trim(), { push: false });

		window.addEventListener("popstate", function () {
			var next = new URLSearchParams(window.location.search);
			selectProject(String(next.get("project") || "").trim(), { push: false });
		});
	};

	document.addEventListener("DOMContentLoaded", function () {
		main();
	});
})();
