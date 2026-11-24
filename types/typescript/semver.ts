import type { SoftwareVersion, } from "./software";

export declare type SemVerFlag = "" | "^" | "~" | ">" | ">=" | "<" | "<=" | "=";

export declare type PreReleaseOrBuild = `-${string}` | `+${string}` | "";

export declare type WildcardVersion = "*" | "latest" | `${number}.x` | `${number}.${number}.x`;

export declare type SemVerVersion = SoftwareVersion & ( `${SemVerFlag}${number}.${number}.${number}${PreReleaseOrBuild}` | WildcardVersion );
