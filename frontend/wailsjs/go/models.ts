export namespace builtin {
	
	export class Http {
	    Disable: boolean;
	    Approval: string;
	    Y: http.CallResult;
	    Z: http.CallArguments;
	
	    static createFrom(source: any = {}) {
	        return new Http(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], http.CallResult);
	        this.Z = this.convertValues(source["Z"], http.CallArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CommandExecution {
	    Disable: boolean;
	    Approval: string;
	    Z: command.CommandExecutionArguments;
	
	    static createFrom(source: any = {}) {
	        return new CommandExecution(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Z = this.convertValues(source["Z"], command.CommandExecutionArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DirectoryDeletion {
	    Disable: boolean;
	    Approval: string;
	    Y: file.DirectoryDeletionResult;
	    Z: file.DirectoryDeletionArguments;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryDeletion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], file.DirectoryDeletionResult);
	        this.Z = this.convertValues(source["Z"], file.DirectoryDeletionArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DirectoryTempCreation {
	    Disable: boolean;
	    Approval: string;
	    Y: file.DirectoryTempCreationResult;
	    // Go type: file
	    Z: any;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryTempCreation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], file.DirectoryTempCreationResult);
	        this.Z = this.convertValues(source["Z"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DirectoryCreation {
	    Disable: boolean;
	    Approval: string;
	    Y: file.DirectoryCreationResult;
	    Z: file.DirectoryCreationArguments;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryCreation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], file.DirectoryCreationResult);
	        this.Z = this.convertValues(source["Z"], file.DirectoryCreationArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FileDeletion {
	    Disable: boolean;
	    Approval: string;
	    Y: file.FileDeletionResult;
	    Z: file.FileDeletionArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileDeletion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], file.FileDeletionResult);
	        this.Z = this.convertValues(source["Z"], file.FileDeletionArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FileReading {
	    Disable: boolean;
	    Approval: string;
	    Y: file.FileReadingResult;
	    Z: file.FileReadingArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileReading(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], file.FileReadingResult);
	        this.Z = this.convertValues(source["Z"], file.FileReadingArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FileAppending {
	    Disable: boolean;
	    Approval: string;
	    Y: file.FileAppendingResult;
	    Z: file.FileAppendingArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileAppending(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], file.FileAppendingResult);
	        this.Z = this.convertValues(source["Z"], file.FileAppendingArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FileTempCreation {
	    Disable: boolean;
	    Approval: string;
	    Y: file.FileTempCreationResult;
	    Z: file.FileTempCreationArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileTempCreation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], file.FileTempCreationResult);
	        this.Z = this.convertValues(source["Z"], file.FileTempCreationArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FileCreation {
	    Disable: boolean;
	    Approval: string;
	    Y: file.FileCreationResult;
	    Z: file.FileCreationArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileCreation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], file.FileCreationResult);
	        this.Z = this.convertValues(source["Z"], file.FileCreationArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ChangeTimes {
	    Disable: boolean;
	    Approval: string;
	    // Go type: file
	    Y: any;
	    Z: file.ChangeTimesArguments;
	
	    static createFrom(source: any = {}) {
	        return new ChangeTimes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], null);
	        this.Z = this.convertValues(source["Z"], file.ChangeTimesArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ChangeOwner {
	    Disable: boolean;
	    Approval: string;
	    // Go type: file
	    Y: any;
	    Z: file.ChangeOwnerArguments;
	
	    static createFrom(source: any = {}) {
	        return new ChangeOwner(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], null);
	        this.Z = this.convertValues(source["Z"], file.ChangeOwnerArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ChangeMode {
	    Disable: boolean;
	    Approval: string;
	    // Go type: file
	    Y: any;
	    Z: file.ChangeModeArguments;
	
	    static createFrom(source: any = {}) {
	        return new ChangeMode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], null);
	        this.Z = this.convertValues(source["Z"], file.ChangeModeArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Stats {
	    Disable: boolean;
	    Approval: string;
	    Y: file.StatsResult;
	    Z: file.StatsArguments;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], file.StatsResult);
	        this.Z = this.convertValues(source["Z"], file.StatsArguments);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SystemTime {
	    Disable: boolean;
	    Approval: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemTime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	    }
	}
	export class Environment {
	    Disable: boolean;
	    Approval: string;
	    Y: system.EnvironmentResult;
	    // Go type: system
	    Z: any;
	
	    static createFrom(source: any = {}) {
	        return new Environment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], system.EnvironmentResult);
	        this.Z = this.convertValues(source["Z"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SystemInfo {
	    Disable: boolean;
	    Approval: string;
	    Y: system.SystemInfoResult;
	    // Go type: system
	    Z: any;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.Approval = source["Approval"];
	        this.Y = this.convertValues(source["Y"], system.SystemInfoResult);
	        this.Z = this.convertValues(source["Z"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BuiltIns {
	    SystemInfo: SystemInfo;
	    Environment: Environment;
	    SystemTime: SystemTime;
	    Stats: Stats;
	    ChangeMode: ChangeMode;
	    ChangeOwner: ChangeOwner;
	    ChangeTimes: ChangeTimes;
	    FileCreation: FileCreation;
	    FileTempCreation: FileTempCreation;
	    FileAppending: FileAppending;
	    FileReading: FileReading;
	    FileDeletion: FileDeletion;
	    DirectoryCreation: DirectoryCreation;
	    DirectoryTempCreation: DirectoryTempCreation;
	    DirectoryDeletion: DirectoryDeletion;
	    CommandExec: CommandExecution;
	    Http: Http;
	
	    static createFrom(source: any = {}) {
	        return new BuiltIns(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SystemInfo = this.convertValues(source["SystemInfo"], SystemInfo);
	        this.Environment = this.convertValues(source["Environment"], Environment);
	        this.SystemTime = this.convertValues(source["SystemTime"], SystemTime);
	        this.Stats = this.convertValues(source["Stats"], Stats);
	        this.ChangeMode = this.convertValues(source["ChangeMode"], ChangeMode);
	        this.ChangeOwner = this.convertValues(source["ChangeOwner"], ChangeOwner);
	        this.ChangeTimes = this.convertValues(source["ChangeTimes"], ChangeTimes);
	        this.FileCreation = this.convertValues(source["FileCreation"], FileCreation);
	        this.FileTempCreation = this.convertValues(source["FileTempCreation"], FileTempCreation);
	        this.FileAppending = this.convertValues(source["FileAppending"], FileAppending);
	        this.FileReading = this.convertValues(source["FileReading"], FileReading);
	        this.FileDeletion = this.convertValues(source["FileDeletion"], FileDeletion);
	        this.DirectoryCreation = this.convertValues(source["DirectoryCreation"], DirectoryCreation);
	        this.DirectoryTempCreation = this.convertValues(source["DirectoryTempCreation"], DirectoryTempCreation);
	        this.DirectoryDeletion = this.convertValues(source["DirectoryDeletion"], DirectoryDeletion);
	        this.CommandExec = this.convertValues(source["CommandExec"], CommandExecution);
	        this.Http = this.convertValues(source["Http"], Http);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

}

export namespace command {
	
	export class CommandExecutionArguments {
	    command: string;
	    working_directory: string;
	    environment: Record<string, string>;
	    out: boolean;
	    err: boolean;
	    first: number;
	    last: number;
	
	    static createFrom(source: any = {}) {
	        return new CommandExecutionArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.command = source["command"];
	        this.working_directory = source["working_directory"];
	        this.environment = source["environment"];
	        this.out = source["out"];
	        this.err = source["err"];
	        this.first = source["first"];
	        this.last = source["last"];
	    }
	}
	export class FunctionDefinition {
	    name: string;
	    description: string;
	    parameters: mcp.ToolInputSchema;
	    approval: string;
	    command?: string;
	    commandExpr?: string;
	    env?: Record<string, string>;
	    additionalEnv?: Record<string, string>;
	    workingDir?: string;
	
	    static createFrom(source: any = {}) {
	        return new FunctionDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.parameters = this.convertValues(source["parameters"], mcp.ToolInputSchema);
	        this.approval = source["approval"];
	        this.command = source["command"];
	        this.commandExpr = source["commandExpr"];
	        this.env = source["env"];
	        this.additionalEnv = source["additionalEnv"];
	        this.workingDir = source["workingDir"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace common {
	
	export class NumberContainer {
	    Expression?: string;
	    Value?: number;
	
	    static createFrom(source: any = {}) {
	        return new NumberContainer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Expression = source["Expression"];
	        this.Value = source["Value"];
	    }
	}
	export class SecretCommand {
	    Name: string;
	    Args: string[];
	    Env: Record<string, string>;
	    NoTrim: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SecretCommand(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Args = source["Args"];
	        this.Env = source["Env"];
	        this.NoTrim = source["NoTrim"];
	    }
	}
	export class Secret {
	    Plain: string;
	    Command: SecretCommand;
	
	    static createFrom(source: any = {}) {
	        return new Secret(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Plain = source["Plain"];
	        this.Command = this.convertValues(source["Command"], SecretCommand);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class StringContainer {
	    Expression?: string;
	    Value?: string;
	
	    static createFrom(source: any = {}) {
	        return new StringContainer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Expression = source["Expression"];
	        this.Value = source["Value"];
	    }
	}

}

export namespace controller {
	
	export class ApplicationConfig {
	    Config: model.Config;
	    ActiveProfile: model.Profile;
	    Z: mcp.Tool;
	
	    static createFrom(source: any = {}) {
	        return new ApplicationConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Config = this.convertValues(source["Config"], model.Config);
	        this.ActiveProfile = this.convertValues(source["ActiveProfile"], model.Profile);
	        this.Z = this.convertValues(source["Z"], mcp.Tool);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AssetMeta {
	    Path: string;
	    Url: string;
	    MimeType: string;
	
	    static createFrom(source: any = {}) {
	        return new AssetMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Url = source["Url"];
	        this.MimeType = source["MimeType"];
	    }
	}
	export class LLMMessageCallResult {
	    Content: string;
	    Error: string;
	    DurationMs: number;
	    V: mcp.CallToolResult;
	    W: mcp.TextContent;
	    X: mcp.ImageContent;
	    Y: mcp.AudioContent;
	    Z: mcp.EmbeddedResource;
	
	    static createFrom(source: any = {}) {
	        return new LLMMessageCallResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Content = source["Content"];
	        this.Error = source["Error"];
	        this.DurationMs = source["DurationMs"];
	        this.V = this.convertValues(source["V"], mcp.CallToolResult);
	        this.W = this.convertValues(source["W"], mcp.TextContent);
	        this.X = this.convertValues(source["X"], mcp.ImageContent);
	        this.Y = this.convertValues(source["Y"], mcp.AudioContent);
	        this.Z = this.convertValues(source["Z"], mcp.EmbeddedResource);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LLMMessageCallMeta {
	    BuiltIn: boolean;
	    Custom: boolean;
	    Mcp: boolean;
	    NeedsApproval: boolean;
	    ToolName: string;
	    ToolDescription: string;
	
	    static createFrom(source: any = {}) {
	        return new LLMMessageCallMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BuiltIn = source["BuiltIn"];
	        this.Custom = source["Custom"];
	        this.Mcp = source["Mcp"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.ToolName = source["ToolName"];
	        this.ToolDescription = source["ToolDescription"];
	    }
	}
	export class LLMMessageCall {
	    Id: string;
	    Function: string;
	    Arguments: string;
	    Meta: LLMMessageCallMeta;
	    Result?: LLMMessageCallResult;
	
	    static createFrom(source: any = {}) {
	        return new LLMMessageCall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Function = source["Function"];
	        this.Arguments = source["Arguments"];
	        this.Meta = this.convertValues(source["Meta"], LLMMessageCallMeta);
	        this.Result = this.convertValues(source["Result"], LLMMessageCallResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LLMMessageContentPart {
	    Type: string;
	    Content: string;
	    Call: LLMMessageCall;
	
	    static createFrom(source: any = {}) {
	        return new LLMMessageContentPart(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.Content = source["Content"];
	        this.Call = this.convertValues(source["Call"], LLMMessageCall);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LLMMessage {
	    Id: string;
	    Role: string;
	    ContentParts: LLMMessageContentPart[];
	    Created: number;
	
	    static createFrom(source: any = {}) {
	        return new LLMMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Role = source["Role"];
	        this.ContentParts = this.convertValues(source["ContentParts"], LLMMessageContentPart);
	        this.Created = source["Created"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LLMAskArgs {
	    History: LLMMessage[];
	
	    static createFrom(source: any = {}) {
	        return new LLMAskArgs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.History = this.convertValues(source["History"], LLMMessage);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LLMAskResult {
	    Content: string;
	    Consumption: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new LLMAskResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Content = source["Content"];
	        this.Consumption = source["Consumption"];
	    }
	}
	
	
	
	
	
	export class OpenFileDialogArgs {
	    Title: string;
	
	    static createFrom(source: any = {}) {
	        return new OpenFileDialogArgs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	    }
	}

}

export namespace file {
	
	export class ChangeModeArguments {
	    path: string;
	    permission: string;
	
	    static createFrom(source: any = {}) {
	        return new ChangeModeArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.permission = source["permission"];
	    }
	}
	export class ChangeOwnerArguments {
	    path: string;
	    user_id?: number;
	    group_id?: number;
	
	    static createFrom(source: any = {}) {
	        return new ChangeOwnerArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.user_id = source["user_id"];
	        this.group_id = source["group_id"];
	    }
	}
	export class ChangeTimesArguments {
	    path: string;
	    access_time: string;
	    modification_time: string;
	
	    static createFrom(source: any = {}) {
	        return new ChangeTimesArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.access_time = source["access_time"];
	        this.modification_time = source["modification_time"];
	    }
	}
	export class DirectoryCreationArguments {
	    path: string;
	    permission: string;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryCreationArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.permission = source["permission"];
	    }
	}
	export class DirectoryCreationResult {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryCreationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class DirectoryDeletionArguments {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryDeletionArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class DirectoryDeletionResult {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryDeletionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class DirectoryTempCreationResult {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryTempCreationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class FileAppendingArguments {
	    path: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new FileAppendingArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.content = source["content"];
	    }
	}
	export class FileAppendingResult {
	    path: string;
	    written: number;
	
	    static createFrom(source: any = {}) {
	        return new FileAppendingResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.written = source["written"];
	    }
	}
	export class FileCreationArguments {
	    path: string;
	    content: string;
	    permission: string;
	
	    static createFrom(source: any = {}) {
	        return new FileCreationArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.content = source["content"];
	        this.permission = source["permission"];
	    }
	}
	export class FileCreationResult {
	    path: string;
	    written: number;
	
	    static createFrom(source: any = {}) {
	        return new FileCreationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.written = source["written"];
	    }
	}
	export class FileDeletionArguments {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new FileDeletionArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class FileDeletionResult {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new FileDeletionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class FileReadingLimits {
	    m: string;
	    o: number;
	    l: number;
	
	    static createFrom(source: any = {}) {
	        return new FileReadingLimits(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.m = source["m"];
	        this.o = source["o"];
	        this.l = source["l"];
	    }
	}
	export class FileReadingArguments {
	    path: string;
	    limits?: FileReadingLimits;
	    lm: string;
	    lo: number;
	    ll: number;
	
	    static createFrom(source: any = {}) {
	        return new FileReadingArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.limits = this.convertValues(source["limits"], FileReadingLimits);
	        this.lm = source["lm"];
	        this.lo = source["lo"];
	        this.ll = source["ll"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class FileReadingResult {
	    path: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new FileReadingResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.content = source["content"];
	    }
	}
	export class FileTempCreationArguments {
	    content: string;
	    suffix: string;
	    permission: string;
	
	    static createFrom(source: any = {}) {
	        return new FileTempCreationArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.suffix = source["suffix"];
	        this.permission = source["permission"];
	    }
	}
	export class FileTempCreationResult {
	    path: string;
	    written: number;
	
	    static createFrom(source: any = {}) {
	        return new FileTempCreationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.written = source["written"];
	    }
	}
	export class StatsArguments {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new StatsArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}
	export class StatsResult {
	    path: string;
	    isDirectory: boolean;
	    isRegular: boolean;
	    permissions: string;
	    size: number;
	    // Go type: time
	    modTime: any;
	
	    static createFrom(source: any = {}) {
	        return new StatsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.isDirectory = source["isDirectory"];
	        this.isRegular = source["isRegular"];
	        this.permissions = source["permissions"];
	        this.size = source["size"];
	        this.modTime = this.convertValues(source["modTime"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace history {
	
	export class MessageCallResult {
	    c: string;
	    e: string;
	    d: number;
	
	    static createFrom(source: any = {}) {
	        return new MessageCallResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.c = source["c"];
	        this.e = source["e"];
	        this.d = source["d"];
	    }
	}
	export class MessageCall {
	    i?: string;
	    f?: string;
	    a?: string;
	    r?: MessageCallResult;
	
	    static createFrom(source: any = {}) {
	        return new MessageCall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.i = source["i"];
	        this.f = source["f"];
	        this.a = source["a"];
	        this.r = this.convertValues(source["r"], MessageCallResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MessageContentPart {
	    t?: string;
	    c?: string;
	    ca?: MessageCall;
	
	    static createFrom(source: any = {}) {
	        return new MessageContentPart(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.t = source["t"];
	        this.c = source["c"];
	        this.ca = this.convertValues(source["ca"], MessageCall);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Message {
	    i?: string;
	    r?: string;
	    p?: MessageContentPart[];
	    t?: number;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.i = source["i"];
	        this.r = source["r"];
	        this.p = this.convertValues(source["p"], MessageContentPart);
	        this.t = source["t"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EntryContent {
	    m: Message[];
	
	    static createFrom(source: any = {}) {
	        return new EntryContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.m = this.convertValues(source["m"], Message);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EntryMeta {
	    v: number;
	    t: number;
	
	    static createFrom(source: any = {}) {
	        return new EntryMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.v = source["v"];
	        this.t = source["t"];
	    }
	}
	export class Entry {
	    m: EntryMeta;
	    c: EntryContent;
	
	    static createFrom(source: any = {}) {
	        return new Entry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.m = this.convertValues(source["m"], EntryMeta);
	        this.c = this.convertValues(source["c"], EntryContent);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	

}

export namespace http {
	
	export class CallArguments {
	    method: string;
	    url: string;
	    header: Record<string, string>;
	    body: string;
	
	    static createFrom(source: any = {}) {
	        return new CallArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.method = source["method"];
	        this.url = source["url"];
	        this.header = source["header"];
	        this.body = source["body"];
	    }
	}
	export class CallResult {
	    status_code: number;
	    status: string;
	    header: Record<string, Array<string>>;
	    body: string;
	
	    static createFrom(source: any = {}) {
	        return new CallResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status_code = source["status_code"];
	        this.status = source["status"];
	        this.header = source["header"];
	        this.body = source["body"];
	    }
	}

}

export namespace llm {
	
	export class AnthropicCache {
	    SystemMessage?: boolean;
	    Tools?: boolean;
	    Chat?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AnthropicCache(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SystemMessage = source["SystemMessage"];
	        this.Tools = source["Tools"];
	        this.Chat = source["Chat"];
	    }
	}
	export class AnthropicConfig {
	    Token: common.Secret;
	    BaseUrl: string;
	    Model: string;
	    Cache: AnthropicCache;
	
	    static createFrom(source: any = {}) {
	        return new AnthropicConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Token = this.convertValues(source["Token"], common.Secret);
	        this.BaseUrl = source["BaseUrl"];
	        this.Model = source["Model"];
	        this.Cache = this.convertValues(source["Cache"], AnthropicCache);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AnythingLLMThreadConfig {
	    Name: common.StringContainer;
	    Delete: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AnythingLLMThreadConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = this.convertValues(source["Name"], common.StringContainer);
	        this.Delete = source["Delete"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AnythingLLMConfig {
	    BaseURL: string;
	    Token: common.Secret;
	    Workspace: string;
	    Thread: AnythingLLMThreadConfig;
	
	    static createFrom(source: any = {}) {
	        return new AnythingLLMConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BaseURL = source["BaseURL"];
	        this.Token = this.convertValues(source["Token"], common.Secret);
	        this.Workspace = source["Workspace"];
	        this.Thread = this.convertValues(source["Thread"], AnythingLLMThreadConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ToolCall {
	    Server: string;
	    Name: string;
	    Arguments: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new ToolCall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Server = source["Server"];
	        this.Name = source["Name"];
	        this.Arguments = source["Arguments"];
	    }
	}
	export class Message {
	    Role: string;
	    Content: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Role = source["Role"];
	        this.Content = source["Content"];
	    }
	}
	export class PromptConfig {
	    System: string;
	    InitMessages: Message[];
	    InitToolCalls: ToolCall[];
	    InitValue: string;
	    InitAttachments: string[];
	
	    static createFrom(source: any = {}) {
	        return new PromptConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.System = source["System"];
	        this.InitMessages = this.convertValues(source["InitMessages"], Message);
	        this.InitToolCalls = this.convertValues(source["InitToolCalls"], ToolCall);
	        this.InitValue = source["InitValue"];
	        this.InitAttachments = source["InitAttachments"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CallOptionsConfig {
	    Prompt: PromptConfig;
	    MaxToken: number;
	    Temperature: number;
	    TopK: number;
	    TopP: number;
	    MinLength: number;
	    MaxLength: number;
	
	    static createFrom(source: any = {}) {
	        return new CallOptionsConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Prompt = this.convertValues(source["Prompt"], PromptConfig);
	        this.MaxToken = source["MaxToken"];
	        this.Temperature = source["Temperature"];
	        this.TopK = source["TopK"];
	        this.TopP = source["TopP"];
	        this.MinLength = source["MinLength"];
	        this.MaxLength = source["MaxLength"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DeepSeekConfig {
	    APIKey: common.Secret;
	    Model: string;
	    BaseUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new DeepSeekConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APIKey = this.convertValues(source["APIKey"], common.Secret);
	        this.Model = source["Model"];
	        this.BaseUrl = source["BaseUrl"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GoogleAIConfig {
	    APIKey: common.Secret;
	    Model: string;
	    HarmThreshold?: number;
	    ToolCacheTTL?: number;
	
	    static createFrom(source: any = {}) {
	        return new GoogleAIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APIKey = this.convertValues(source["APIKey"], common.Secret);
	        this.Model = source["Model"];
	        this.HarmThreshold = source["HarmThreshold"];
	        this.ToolCacheTTL = source["ToolCacheTTL"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MistralConfig {
	    ApiKey: common.Secret;
	    Endpoint: string;
	    Model: string;
	
	    static createFrom(source: any = {}) {
	        return new MistralConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ApiKey = this.convertValues(source["ApiKey"], common.Secret);
	        this.Endpoint = source["Endpoint"];
	        this.Model = source["Model"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OllamaConfig {
	    ServerURL: string;
	    Model: string;
	
	    static createFrom(source: any = {}) {
	        return new OllamaConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ServerURL = source["ServerURL"];
	        this.Model = source["Model"];
	    }
	}
	export class OpenAIConfig {
	    APIKey: common.Secret;
	    APIType: string;
	    APIVersion: string;
	    Model: string;
	    BaseUrl: string;
	    Organization: string;
	
	    static createFrom(source: any = {}) {
	        return new OpenAIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APIKey = this.convertValues(source["APIKey"], common.Secret);
	        this.APIType = source["APIType"];
	        this.APIVersion = source["APIVersion"];
	        this.Model = source["Model"];
	        this.BaseUrl = source["BaseUrl"];
	        this.Organization = source["Organization"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LocalAIConfig {
	    APIKey: common.Secret;
	    Model: string;
	    BaseUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new LocalAIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APIKey = this.convertValues(source["APIKey"], common.Secret);
	        this.Model = source["Model"];
	        this.BaseUrl = source["BaseUrl"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CopilotConfig {
	
	
	    static createFrom(source: any = {}) {
	        return new CopilotConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class LLMConfig {
	    Backend: string;
	    // Go type: CopilotConfig
	    Copilot: any;
	    LocalAI: LocalAIConfig;
	    OpenAI: OpenAIConfig;
	    AnythingLLM: AnythingLLMConfig;
	    Ollama: OllamaConfig;
	    Mistral: MistralConfig;
	    Anthropic: AnthropicConfig;
	    DeepSeek: DeepSeekConfig;
	    Google: GoogleAIConfig;
	    CallOptions: CallOptionsConfig;
	    Tool: tools.Config;
	
	    static createFrom(source: any = {}) {
	        return new LLMConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Backend = source["Backend"];
	        this.Copilot = this.convertValues(source["Copilot"], null);
	        this.LocalAI = this.convertValues(source["LocalAI"], LocalAIConfig);
	        this.OpenAI = this.convertValues(source["OpenAI"], OpenAIConfig);
	        this.AnythingLLM = this.convertValues(source["AnythingLLM"], AnythingLLMConfig);
	        this.Ollama = this.convertValues(source["Ollama"], OllamaConfig);
	        this.Mistral = this.convertValues(source["Mistral"], MistralConfig);
	        this.Anthropic = this.convertValues(source["Anthropic"], AnthropicConfig);
	        this.DeepSeek = this.convertValues(source["DeepSeek"], DeepSeekConfig);
	        this.Google = this.convertValues(source["Google"], GoogleAIConfig);
	        this.CallOptions = this.convertValues(source["CallOptions"], CallOptionsConfig);
	        this.Tool = this.convertValues(source["Tool"], tools.Config);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	

}

export namespace mcp {
	
	export class Annotations {
	    audience?: string[];
	    priority?: number;
	
	    static createFrom(source: any = {}) {
	        return new Annotations(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.audience = source["audience"];
	        this.priority = source["priority"];
	    }
	}
	export class AudioContent {
	    // Go type: Annotations
	    annotations?: any;
	    type: string;
	    data: string;
	    mimeType: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.annotations = this.convertValues(source["annotations"], null);
	        this.type = source["type"];
	        this.data = source["data"];
	        this.mimeType = source["mimeType"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CallToolResult {
	    _meta?: Record<string, any>;
	    content: any[];
	    isError?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CallToolResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this._meta = source["_meta"];
	        this.content = source["content"];
	        this.isError = source["isError"];
	    }
	}
	export class EmbeddedResource {
	    // Go type: Annotations
	    annotations?: any;
	    type: string;
	    resource: any;
	
	    static createFrom(source: any = {}) {
	        return new EmbeddedResource(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.annotations = this.convertValues(source["annotations"], null);
	        this.type = source["type"];
	        this.resource = source["resource"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImageContent {
	    // Go type: Annotations
	    annotations?: any;
	    type: string;
	    data: string;
	    mimeType: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.annotations = this.convertValues(source["annotations"], null);
	        this.type = source["type"];
	        this.data = source["data"];
	        this.mimeType = source["mimeType"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Timeout {
	    Init?: number;
	    List?: number;
	    Execution?: number;
	
	    static createFrom(source: any = {}) {
	        return new Timeout(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Init = source["Init"];
	        this.List = source["List"];
	        this.Execution = source["Execution"];
	    }
	}
	export class Server {
	    BaseUrl: string;
	    Headers: Record<string, string>;
	    Name: string;
	    Arguments: string[];
	    Environment: Record<string, string>;
	    Approval: string;
	    Exclude: string[];
	    Timeout: Timeout;
	
	    static createFrom(source: any = {}) {
	        return new Server(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BaseUrl = source["BaseUrl"];
	        this.Headers = source["Headers"];
	        this.Name = source["Name"];
	        this.Arguments = source["Arguments"];
	        this.Environment = source["Environment"];
	        this.Approval = source["Approval"];
	        this.Exclude = source["Exclude"];
	        this.Timeout = this.convertValues(source["Timeout"], Timeout);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TextContent {
	    // Go type: Annotations
	    annotations?: any;
	    type: string;
	    text: string;
	
	    static createFrom(source: any = {}) {
	        return new TextContent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.annotations = this.convertValues(source["annotations"], null);
	        this.type = source["type"];
	        this.text = source["text"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ToolAnnotation {
	    title?: string;
	    readOnlyHint?: boolean;
	    destructiveHint?: boolean;
	    idempotentHint?: boolean;
	    openWorldHint?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ToolAnnotation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.readOnlyHint = source["readOnlyHint"];
	        this.destructiveHint = source["destructiveHint"];
	        this.idempotentHint = source["idempotentHint"];
	        this.openWorldHint = source["openWorldHint"];
	    }
	}
	export class ToolInputSchema {
	    type: string;
	    properties?: Record<string, any>;
	    required?: string[];
	
	    static createFrom(source: any = {}) {
	        return new ToolInputSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.properties = source["properties"];
	        this.required = source["required"];
	    }
	}
	export class Tool {
	    name: string;
	    description?: string;
	    inputSchema: ToolInputSchema;
	    annotations: ToolAnnotation;
	
	    static createFrom(source: any = {}) {
	        return new Tool(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.inputSchema = this.convertValues(source["inputSchema"], ToolInputSchema);
	        this.annotations = this.convertValues(source["annotations"], ToolAnnotation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace model {
	
	export class Help {
	    Arg: boolean;
	    Env: boolean;
	    Yaml: boolean;
	    Styles: boolean;
	    Expr: boolean;
	    Tool: boolean;
	    GenYaml: boolean;
	    DumpYaml: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Help(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Arg = source["Arg"];
	        this.Env = source["Env"];
	        this.Yaml = source["Yaml"];
	        this.Styles = source["Styles"];
	        this.Expr = source["Expr"];
	        this.Tool = source["Tool"];
	        this.GenYaml = source["GenYaml"];
	        this.DumpYaml = source["DumpYaml"];
	    }
	}
	export class Theme {
	    dark: boolean;
	    colors: Record<string, string>;
	    variables: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new Theme(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dark = source["dark"];
	        this.colors = source["colors"];
	        this.variables = source["variables"];
	    }
	}
	export class Themes {
	    dark?: Theme;
	    light?: Theme;
	    custom?: Record<string, Theme>;
	
	    static createFrom(source: any = {}) {
	        return new Themes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dark = this.convertValues(source["dark"], Theme);
	        this.light = this.convertValues(source["light"], Theme);
	        this.custom = this.convertValues(source["custom"], Theme, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class WebKitInspectorConfig {
	    OpenInspectorOnStartup: boolean;
	    HttpServerAddress: string;
	
	    static createFrom(source: any = {}) {
	        return new WebKitInspectorConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OpenInspectorOnStartup = source["OpenInspectorOnStartup"];
	        this.HttpServerAddress = source["HttpServerAddress"];
	    }
	}
	export class VueDevToolsConfig {
	    Host: string;
	    Port: number;
	
	    static createFrom(source: any = {}) {
	        return new VueDevToolsConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Host = source["Host"];
	        this.Port = source["Port"];
	    }
	}
	export class DebugConfig {
	    LogLevel: string;
	    LogLevelParsed?: number;
	    PprofAddress: string;
	    VueDevTools: VueDevToolsConfig;
	    WebKit: WebKitInspectorConfig;
	    DisableCrashDetection: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DebugConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LogLevel = source["LogLevel"];
	        this.LogLevelParsed = source["LogLevelParsed"];
	        this.PprofAddress = source["PprofAddress"];
	        this.VueDevTools = this.convertValues(source["VueDevTools"], VueDevToolsConfig);
	        this.WebKit = this.convertValues(source["WebKit"], WebKitInspectorConfig);
	        this.DisableCrashDetection = source["DisableCrashDetection"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class History {
	    Path?: string;
	
	    static createFrom(source: any = {}) {
	        return new History(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	    }
	}
	export class FileDialogConfig {
	    DefaultDirectory: string;
	    ShowHiddenFiles?: boolean;
	    CanCreateDirectories?: boolean;
	    ResolveAliases?: boolean;
	    TreatPackagesAsDirectories?: boolean;
	    FilterDisplay: string[];
	    FilterPattern: string[];
	
	    static createFrom(source: any = {}) {
	        return new FileDialogConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DefaultDirectory = source["DefaultDirectory"];
	        this.ShowHiddenFiles = source["ShowHiddenFiles"];
	        this.CanCreateDirectories = source["CanCreateDirectories"];
	        this.ResolveAliases = source["ResolveAliases"];
	        this.TreatPackagesAsDirectories = source["TreatPackagesAsDirectories"];
	        this.FilterDisplay = source["FilterDisplay"];
	        this.FilterPattern = source["FilterPattern"];
	    }
	}
	export class Shortcut {
	    Binding: string[];
	    Code: string[];
	    Alt: boolean[];
	    Ctrl: boolean[];
	    Meta: boolean[];
	    Shift: boolean[];
	
	    static createFrom(source: any = {}) {
	        return new Shortcut(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Binding = source["Binding"];
	        this.Code = source["Code"];
	        this.Alt = source["Alt"];
	        this.Ctrl = source["Ctrl"];
	        this.Meta = source["Meta"];
	        this.Shift = source["Shift"];
	    }
	}
	export class PromptConfig {
	    MinRows?: number;
	    MaxRows?: number;
	    PinTop?: boolean;
	    SubmitShortcut: Shortcut;
	
	    static createFrom(source: any = {}) {
	        return new PromptConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MinRows = source["MinRows"];
	        this.MaxRows = source["MaxRows"];
	        this.PinTop = source["PinTop"];
	        this.SubmitShortcut = this.convertValues(source["SubmitShortcut"], Shortcut);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class WindowBackgroundColor {
	    R?: number;
	    G?: number;
	    B?: number;
	    A?: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowBackgroundColor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.R = source["R"];
	        this.G = source["G"];
	        this.B = source["B"];
	        this.A = source["A"];
	    }
	}
	export class WindowConfig {
	    Title: string;
	    InitialWidth: common.NumberContainer;
	    MaxHeight: common.NumberContainer;
	    InitialPositionX: common.NumberContainer;
	    InitialPositionY: common.NumberContainer;
	    InitialZoom: common.NumberContainer;
	    BackgroundColor: WindowBackgroundColor;
	    StartState?: number;
	    AlwaysOnTop?: boolean;
	    ShowTitleBar?: boolean;
	    TitleBarHeight?: number;
	    GrowTop?: boolean;
	    Frameless?: boolean;
	    Resizeable?: boolean;
	    Translucent: string;
	
	    static createFrom(source: any = {}) {
	        return new WindowConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.InitialWidth = this.convertValues(source["InitialWidth"], common.NumberContainer);
	        this.MaxHeight = this.convertValues(source["MaxHeight"], common.NumberContainer);
	        this.InitialPositionX = this.convertValues(source["InitialPositionX"], common.NumberContainer);
	        this.InitialPositionY = this.convertValues(source["InitialPositionY"], common.NumberContainer);
	        this.InitialZoom = this.convertValues(source["InitialZoom"], common.NumberContainer);
	        this.BackgroundColor = this.convertValues(source["BackgroundColor"], WindowBackgroundColor);
	        this.StartState = source["StartState"];
	        this.AlwaysOnTop = source["AlwaysOnTop"];
	        this.ShowTitleBar = source["ShowTitleBar"];
	        this.TitleBarHeight = source["TitleBarHeight"];
	        this.GrowTop = source["GrowTop"];
	        this.Frameless = source["Frameless"];
	        this.Resizeable = source["Resizeable"];
	        this.Translucent = source["Translucent"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UIConfig {
	    Window: WindowConfig;
	    Prompt: PromptConfig;
	    FileDialog: FileDialogConfig;
	    Stream?: boolean;
	    QuitShortcut: Shortcut;
	    Theme: string;
	    CodeStyle: string;
	    Language: string;
	
	    static createFrom(source: any = {}) {
	        return new UIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Window = this.convertValues(source["Window"], WindowConfig);
	        this.Prompt = this.convertValues(source["Prompt"], PromptConfig);
	        this.FileDialog = this.convertValues(source["FileDialog"], FileDialogConfig);
	        this.Stream = source["Stream"];
	        this.QuitShortcut = this.convertValues(source["QuitShortcut"], Shortcut);
	        this.Theme = source["Theme"];
	        this.CodeStyle = source["CodeStyle"];
	        this.Language = source["Language"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PrinterConfig {
	    Targets: any[];
	    TargetsRaw: string[];
	    Format: string;
	
	    static createFrom(source: any = {}) {
	        return new PrinterConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Targets = source["Targets"];
	        this.TargetsRaw = source["TargetsRaw"];
	        this.Format = source["Format"];
	    }
	}
	export class ProfileMeta {
	    Icon: string;
	    Description: string;
	
	    static createFrom(source: any = {}) {
	        return new ProfileMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Icon = source["Icon"];
	        this.Description = source["Description"];
	    }
	}
	export class Profile {
	    Meta: ProfileMeta;
	    Printer: PrinterConfig;
	    LLM: llm.LLMConfig;
	    UI: UIConfig;
	    History: History;
	    RestartShortcut: Shortcut;
	
	    static createFrom(source: any = {}) {
	        return new Profile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Meta = this.convertValues(source["Meta"], ProfileMeta);
	        this.Printer = this.convertValues(source["Printer"], PrinterConfig);
	        this.LLM = this.convertValues(source["LLM"], llm.LLMConfig);
	        this.UI = this.convertValues(source["UI"], UIConfig);
	        this.History = this.convertValues(source["History"], History);
	        this.RestartShortcut = this.convertValues(source["RestartShortcut"], Shortcut);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Config {
	    Path: string;
	    MainProfile: Profile;
	    DebugConfig: DebugConfig;
	    ActiveProfile: string;
	    Profiles: Record<string, Profile>;
	    Themes: Themes;
	    Version: boolean;
	    Help: Help;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.MainProfile = this.convertValues(source["MainProfile"], Profile);
	        this.DebugConfig = this.convertValues(source["DebugConfig"], DebugConfig);
	        this.ActiveProfile = source["ActiveProfile"];
	        this.Profiles = this.convertValues(source["Profiles"], Profile, true);
	        this.Themes = this.convertValues(source["Themes"], Themes);
	        this.Version = source["Version"];
	        this.Help = this.convertValues(source["Help"], Help);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

}

export namespace system {
	
	export class EnvironmentResult {
	    env: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.env = source["env"];
	    }
	}
	export class SystemInfoResult {
	    os: string;
	    os_info?: any;
	    arch: string;
	    cpus: number;
	    hostname: string;
	    user_dir: string;
	    user_id: number;
	    group_id: number;
	    working_directory: string;
	    process_id: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfoResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.os_info = source["os_info"];
	        this.arch = source["arch"];
	        this.cpus = source["cpus"];
	        this.hostname = source["hostname"];
	        this.user_dir = source["user_dir"];
	        this.user_id = source["user_id"];
	        this.group_id = source["group_id"];
	        this.working_directory = source["working_directory"];
	        this.process_id = source["process_id"];
	    }
	}

}

export namespace tools {
	
	export class Config {
	    BuiltIns: builtin.BuiltIns;
	    McpServer: Record<string, mcp.Server>;
	    Custom: Record<string, command.FunctionDefinition>;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BuiltIns = this.convertValues(source["BuiltIns"], builtin.BuiltIns);
	        this.McpServer = this.convertValues(source["McpServer"], mcp.Server, true);
	        this.Custom = this.convertValues(source["Custom"], command.FunctionDefinition, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

