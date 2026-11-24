export declare type SemVerFlag = "" | "^" | "~" | ">" | ">=" | "<" | "<=" | "=";
export declare type PreReleaseOrBuild = `-${string}` | `+${string}` | "";
export declare type WildcardVersion = "*" | "latest" | `${number}.x` | `${number}.${number}.x`;
export declare type SemVerVersion = `${SemVerFlag}${number}.${number}.${number}${PreReleaseOrBuild}` | WildcardVersion;
