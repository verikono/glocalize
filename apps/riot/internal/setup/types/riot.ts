import type { RelativeDirectoryPath, } from "./filesystem";
import type { SemVerVersion, } from "./semver";



/** supported languages/environments by riot */
export declare type RiotSupportedLanguage = "rust" | "go" | "typescript" | "react" | "lynx" | "python";


/** stardard location names */
export declare type RepoLocation = "apps" | "packages" | "types";


/** a relative directory path anchored to the repository base directory by using a caret(^) to indicate as such */
export declare type RepoRelativeDirectoryPath = `^${RelativeDirectoryPath}`;


/** a location map. */
export declare type RepoLocationMap = Record<Partial<RepoLocation>, RepoRelativeDirectoryPath>;


/**
 * @property {boolean?} active flagging this false disables typing for this language.
 */
export declare type RiotBaseTypingConfig = {
    active: boolean; // true by default but flaggable to turn off.
}


/**
 * @property {boolean?} pydantic flagging this true enables pydantic support
 * @property {boolean?} pyrefly flagging this true enabled pyrefly support
 */
export declare type RiotPythonTypingConfig = RiotBaseTypingConfig & {
    pydantic: boolean; // off by default
    pyrefly: boolean; // off by default
}

export declare type RiotTypingConfig = RiotBaseTypingConfig | RiotPythonTypingConfig;

export declare type RiotLanguageSpec = {
    version: SemVerVersion;
    workspace_file: RepoRelativeDirectoryPath;
    locations: RepoLocationMap;
    typing?: RiotTypingConfig;
}
