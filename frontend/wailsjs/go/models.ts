export namespace config {
	
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
	    LogLevel: number;
	    PprofAddress: string;
	    VueDevTools: VueDevToolsConfig;
	    WebKit: WebKitInspectorConfig;
	    DisableCrashDetection: boolean;
	    RestartShortcut: Shortcut;
	    PrintVersion: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DebugConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LogLevel = source["LogLevel"];
	        this.PprofAddress = source["PprofAddress"];
	        this.VueDevTools = this.convertValues(source["VueDevTools"], VueDevToolsConfig);
	        this.WebKit = this.convertValues(source["WebKit"], WebKitInspectorConfig);
	        this.DisableCrashDetection = source["DisableCrashDetection"];
	        this.RestartShortcut = this.convertValues(source["RestartShortcut"], Shortcut);
	        this.PrintVersion = source["PrintVersion"];
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
	    Format: string;
	    Targets: any[];
	    TargetsRaw: string[];
	
	    static createFrom(source: any = {}) {
	        return new PrinterConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Format = source["Format"];
	        this.Targets = source["Targets"];
	        this.TargetsRaw = source["TargetsRaw"];
	    }
	}
	export class FileDialogConfig {
	    DefaultDirectory: string;
	    ShowHiddenFiles: boolean;
	    CanCreateDirectories: boolean;
	    ResolveAliases: boolean;
	    TreatPackagesAsDirectories: boolean;
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
	    InitValue: string;
	    InitAttachments: string[];
	    MinRows: number;
	    MaxRows: number;
	    PinTop: boolean;
	    SubmitShortcut: Shortcut;
	
	    static createFrom(source: any = {}) {
	        return new PromptConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.InitValue = source["InitValue"];
	        this.InitAttachments = source["InitAttachments"];
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
	    R: number;
	    G: number;
	    B: number;
	    A: number;
	
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
	    InitialWidth: expression.NumberContainer;
	    MaxHeight: expression.NumberContainer;
	    InitialPositionX: expression.NumberContainer;
	    InitialPositionY: expression.NumberContainer;
	    InitialZoom: expression.NumberContainer;
	    BackgroundColor: WindowBackgroundColor;
	    StartState: number;
	    AlwaysOnTop: boolean;
	    GrowTop: boolean;
	    Frameless: boolean;
	    Resizeable: boolean;
	    Translucent: string;
	
	    static createFrom(source: any = {}) {
	        return new WindowConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.InitialWidth = this.convertValues(source["InitialWidth"], expression.NumberContainer);
	        this.MaxHeight = this.convertValues(source["MaxHeight"], expression.NumberContainer);
	        this.InitialPositionX = this.convertValues(source["InitialPositionX"], expression.NumberContainer);
	        this.InitialPositionY = this.convertValues(source["InitialPositionY"], expression.NumberContainer);
	        this.InitialZoom = this.convertValues(source["InitialZoom"], expression.NumberContainer);
	        this.BackgroundColor = this.convertValues(source["BackgroundColor"], WindowBackgroundColor);
	        this.StartState = source["StartState"];
	        this.AlwaysOnTop = source["AlwaysOnTop"];
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
	    Stream: boolean;
	    QuitShortcut: Shortcut;
	    Theme: string;
	    MinMaxPosition: string;
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
	        this.MinMaxPosition = source["MinMaxPosition"];
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
	export class Config {
	    UI: UIConfig;
	    LLM: llm.LLMConfig;
	    Printer: PrinterConfig;
	    Debug: DebugConfig;
	    Config: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UI = this.convertValues(source["UI"], UIConfig);
	        this.LLM = this.convertValues(source["LLM"], llm.LLMConfig);
	        this.Printer = this.convertValues(source["Printer"], PrinterConfig);
	        this.Debug = this.convertValues(source["Debug"], DebugConfig);
	        this.Config = source["Config"];
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

export namespace controller {
	
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
	
	    static createFrom(source: any = {}) {
	        return new LLMMessageCallResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Content = source["Content"];
	        this.Error = source["Error"];
	        this.DurationMs = source["DurationMs"];
	    }
	}
	export class LLMMessageCall {
	    Id: string;
	    Function: string;
	    Arguments: string;
	    NeedsApproval: boolean;
	    BuiltIn: boolean;
	    Result?: LLMMessageCallResult;
	
	    static createFrom(source: any = {}) {
	        return new LLMMessageCall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Function = source["Function"];
	        this.Arguments = source["Arguments"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.BuiltIn = source["BuiltIn"];
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
	
	    static createFrom(source: any = {}) {
	        return new LLMMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Id = source["Id"];
	        this.Role = source["Role"];
	        this.ContentParts = this.convertValues(source["ContentParts"], LLMMessageContentPart);
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

export namespace expression {
	
	export class NumberContainer {
	    Expression: string;
	    Value: number;
	
	    static createFrom(source: any = {}) {
	        return new NumberContainer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Expression = source["Expression"];
	        this.Value = source["Value"];
	    }
	}
	export class StringContainer {
	    Expression: string;
	    Value: string;
	
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

export namespace llm {
	
	export class AnthropicConfig {
	    Token: string;
	    BaseUrl: string;
	    Model: string;
	
	    static createFrom(source: any = {}) {
	        return new AnthropicConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Token = source["Token"];
	        this.BaseUrl = source["BaseUrl"];
	        this.Model = source["Model"];
	    }
	}
	export class AnythingLLMThreadConfig {
	    Name: expression.StringContainer;
	    Delete: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AnythingLLMThreadConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = this.convertValues(source["Name"], expression.StringContainer);
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
	    Token: string;
	    Workspace: string;
	    Thread: AnythingLLMThreadConfig;
	
	    static createFrom(source: any = {}) {
	        return new AnythingLLMConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BaseURL = source["BaseURL"];
	        this.Token = source["Token"];
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
	
	export class CallOptionsConfig {
	    SystemPrompt: string;
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
	        this.SystemPrompt = source["SystemPrompt"];
	        this.MaxToken = source["MaxToken"];
	        this.Temperature = source["Temperature"];
	        this.TopK = source["TopK"];
	        this.TopP = source["TopP"];
	        this.MinLength = source["MinLength"];
	        this.MaxLength = source["MaxLength"];
	    }
	}
	export class DeepSeekConfig {
	    APIKey: string;
	    Model: string;
	    BaseUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new DeepSeekConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APIKey = source["APIKey"];
	        this.Model = source["Model"];
	        this.BaseUrl = source["BaseUrl"];
	    }
	}
	export class MistralConfig {
	    ApiKey: string;
	    Endpoint: string;
	    Model: string;
	
	    static createFrom(source: any = {}) {
	        return new MistralConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ApiKey = source["ApiKey"];
	        this.Endpoint = source["Endpoint"];
	        this.Model = source["Model"];
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
	    APIKey: string;
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
	        this.APIKey = source["APIKey"];
	        this.APIType = source["APIType"];
	        this.APIVersion = source["APIVersion"];
	        this.Model = source["Model"];
	        this.BaseUrl = source["BaseUrl"];
	        this.Organization = source["Organization"];
	    }
	}
	export class LocalAIConfig {
	    APIKey: string;
	    Model: string;
	    BaseUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new LocalAIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APIKey = source["APIKey"];
	        this.Model = source["Model"];
	        this.BaseUrl = source["BaseUrl"];
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
	    CallOptions: CallOptionsConfig;
	    Tools: tools.Config;
	
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
	        this.CallOptions = this.convertValues(source["CallOptions"], CallOptionsConfig);
	        this.Tools = this.convertValues(source["Tools"], tools.Config);
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

export namespace tools {
	
	export class CommandExecutionArguments {
	    name: string;
	    arguments: string[];
	    working_directory: string;
	    environment: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new CommandExecutionArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.arguments = source["arguments"];
	        this.working_directory = source["working_directory"];
	        this.environment = source["environment"];
	    }
	}
	export class CommandExecution {
	    Disable: boolean;
	    NoApproval: boolean;
	    NoApprovalCommands: string[];
	    NoApprovalCommandsExpr: string[];
	    Z: CommandExecutionArguments;
	
	    static createFrom(source: any = {}) {
	        return new CommandExecution(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NoApproval = source["NoApproval"];
	        this.NoApprovalCommands = source["NoApprovalCommands"];
	        this.NoApprovalCommandsExpr = source["NoApprovalCommandsExpr"];
	        this.Z = this.convertValues(source["Z"], CommandExecutionArguments);
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
	export class DirectoryDeletion {
	    Disable: boolean;
	    NoApproval: boolean;
	    Y: DirectoryDeletionResult;
	    Z: DirectoryDeletionArguments;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryDeletion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NoApproval = source["NoApproval"];
	        this.Y = this.convertValues(source["Y"], DirectoryDeletionResult);
	        this.Z = this.convertValues(source["Z"], DirectoryDeletionArguments);
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
	export class DirectoryTempCreationArguments {
	
	
	    static createFrom(source: any = {}) {
	        return new DirectoryTempCreationArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
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
	export class DirectoryTempCreation {
	    Disable: boolean;
	    NeedsApproval: boolean;
	    Y: DirectoryTempCreationResult;
	    // Go type: DirectoryTempCreationArguments
	    Z: any;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryTempCreation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.Y = this.convertValues(source["Y"], DirectoryTempCreationResult);
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
	export class DirectoryCreation {
	    Disable: boolean;
	    NeedsApproval: boolean;
	    Y: DirectoryCreationResult;
	    Z: DirectoryCreationArguments;
	
	    static createFrom(source: any = {}) {
	        return new DirectoryCreation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.Y = this.convertValues(source["Y"], DirectoryCreationResult);
	        this.Z = this.convertValues(source["Z"], DirectoryCreationArguments);
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
	export class FileDeletion {
	    Disable: boolean;
	    NoApproval: boolean;
	    Y: FileDeletionResult;
	    Z: FileDeletionArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileDeletion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NoApproval = source["NoApproval"];
	        this.Y = this.convertValues(source["Y"], FileDeletionResult);
	        this.Z = this.convertValues(source["Z"], FileDeletionArguments);
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
	
	    static createFrom(source: any = {}) {
	        return new FileReadingArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.limits = this.convertValues(source["limits"], FileReadingLimits);
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
	export class FileReading {
	    Disable: boolean;
	    NeedsApproval: boolean;
	    Y: FileReadingResult;
	    Z: FileReadingArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileReading(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.Y = this.convertValues(source["Y"], FileReadingResult);
	        this.Z = this.convertValues(source["Z"], FileReadingArguments);
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
	export class FileAppending {
	    Disable: boolean;
	    NoApproval: boolean;
	    Y: FileAppendingResult;
	    Z: FileAppendingArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileAppending(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NoApproval = source["NoApproval"];
	        this.Y = this.convertValues(source["Y"], FileAppendingResult);
	        this.Z = this.convertValues(source["Z"], FileAppendingArguments);
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
	export class FileTempCreation {
	    Disable: boolean;
	    NeedsApproval: boolean;
	    Y: FileTempCreationResult;
	    Z: FileTempCreationArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileTempCreation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.Y = this.convertValues(source["Y"], FileTempCreationResult);
	        this.Z = this.convertValues(source["Z"], FileTempCreationArguments);
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
	export class FileCreation {
	    Disable: boolean;
	    NeedsApproval: boolean;
	    Y: FileCreationResult;
	    Z: FileCreationArguments;
	
	    static createFrom(source: any = {}) {
	        return new FileCreation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.Y = this.convertValues(source["Y"], FileCreationResult);
	        this.Z = this.convertValues(source["Z"], FileCreationArguments);
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
	export class ChangeTimesResult {
	
	
	    static createFrom(source: any = {}) {
	        return new ChangeTimesResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class ChangeTimes {
	    Disable: boolean;
	    NeedsApproval: boolean;
	    // Go type: ChangeTimesResult
	    Y: any;
	    Z: ChangeTimesArguments;
	
	    static createFrom(source: any = {}) {
	        return new ChangeTimes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.Y = this.convertValues(source["Y"], null);
	        this.Z = this.convertValues(source["Z"], ChangeTimesArguments);
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
	export class ChangeOwnerArguments {
	    path: string;
	    user_id: number;
	    group_id: number;
	
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
	export class ChangeOwnerResult {
	
	
	    static createFrom(source: any = {}) {
	        return new ChangeOwnerResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class ChangeOwner {
	    Disable: boolean;
	    NoApproval: boolean;
	    // Go type: ChangeOwnerResult
	    Y: any;
	    Z: ChangeOwnerArguments;
	
	    static createFrom(source: any = {}) {
	        return new ChangeOwner(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NoApproval = source["NoApproval"];
	        this.Y = this.convertValues(source["Y"], null);
	        this.Z = this.convertValues(source["Z"], ChangeOwnerArguments);
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
	export class ChangeModeResult {
	
	
	    static createFrom(source: any = {}) {
	        return new ChangeModeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class ChangeMode {
	    Disable: boolean;
	    NoApproval: boolean;
	    // Go type: ChangeModeResult
	    Y: any;
	    Z: ChangeModeArguments;
	
	    static createFrom(source: any = {}) {
	        return new ChangeMode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NoApproval = source["NoApproval"];
	        this.Y = this.convertValues(source["Y"], null);
	        this.Z = this.convertValues(source["Z"], ChangeModeArguments);
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
	export class Stats {
	    Disable: boolean;
	    NeedsApproval: boolean;
	    Y: StatsResult;
	    Z: StatsArguments;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.Y = this.convertValues(source["Y"], StatsResult);
	        this.Z = this.convertValues(source["Z"], StatsArguments);
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
	    NeedsApproval: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SystemTime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	    }
	}
	export class EnvironmentArguments {
	
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
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
	export class Environment {
	    Disable: boolean;
	    NeedsApproval: boolean;
	    Y: EnvironmentResult;
	    // Go type: EnvironmentArguments
	    Z: any;
	
	    static createFrom(source: any = {}) {
	        return new Environment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.Y = this.convertValues(source["Y"], EnvironmentResult);
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
	export class SystemInfoArguments {
	
	
	    static createFrom(source: any = {}) {
	        return new SystemInfoArguments(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class SystemInfoResult {
	    os: string;
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
	export class SystemInfo {
	    Disable: boolean;
	    NeedsApproval: boolean;
	    Y: SystemInfoResult;
	    // Go type: SystemInfoArguments
	    Z: any;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Disable = source["Disable"];
	        this.NeedsApproval = source["NeedsApproval"];
	        this.Y = this.convertValues(source["Y"], SystemInfoResult);
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
	    Disable: boolean;
	
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
	        this.Disable = source["Disable"];
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
	
	
	
	
	
	
	
	
	export class FunctionDefinition {
	    name: string;
	    description: string;
	    parameters: any;
	    approval: boolean;
	    command: string;
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
	        this.parameters = source["parameters"];
	        this.approval = source["approval"];
	        this.command = source["command"];
	        this.env = source["env"];
	        this.additionalEnv = source["additionalEnv"];
	        this.workingDir = source["workingDir"];
	    }
	}
	export class Config {
	    RawTools: string[];
	    Tools: Record<string, FunctionDefinition>;
	    BuiltInTools: BuiltIns;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RawTools = source["RawTools"];
	        this.Tools = this.convertValues(source["Tools"], FunctionDefinition, true);
	        this.BuiltInTools = this.convertValues(source["BuiltInTools"], BuiltIns);
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

