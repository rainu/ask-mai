export namespace config {
	
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
	export class AnythingLLMConfig {
	    BaseURL: string;
	    Token: string;
	    Workspace: string;
	
	    static createFrom(source: any = {}) {
	        return new AnythingLLMConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.BaseURL = source["BaseURL"];
	        this.Token = source["Token"];
	        this.Workspace = source["Workspace"];
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
	export class PrinterConfig {
	    Format: string;
	    Targets: any[];
	
	    static createFrom(source: any = {}) {
	        return new PrinterConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Format = source["Format"];
	        this.Targets = source["Targets"];
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
	export class Shortcut {
	    Code: string;
	    Alt: boolean;
	    Ctrl: boolean;
	    Meta: boolean;
	    Shift: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Shortcut(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Code = source["Code"];
	        this.Alt = source["Alt"];
	        this.Ctrl = source["Ctrl"];
	        this.Meta = source["Meta"];
	        this.Shift = source["Shift"];
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
	    InitialWidth: string;
	    MaxHeight: string;
	    InitialPositionX: string;
	    InitialPositionY: string;
	    InitialZoom: number;
	    BackgroundColor: WindowBackgroundColor;
	    StartState: number;
	    AlwaysOnTop: boolean;
	    Frameless: boolean;
	    Resizeable: boolean;
	    Translucent: string;
	
	    static createFrom(source: any = {}) {
	        return new WindowConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.InitialWidth = source["InitialWidth"];
	        this.MaxHeight = source["MaxHeight"];
	        this.InitialPositionX = source["InitialPositionX"];
	        this.InitialPositionY = source["InitialPositionY"];
	        this.InitialZoom = source["InitialZoom"];
	        this.BackgroundColor = this.convertValues(source["BackgroundColor"], WindowBackgroundColor);
	        this.StartState = source["StartState"];
	        this.AlwaysOnTop = source["AlwaysOnTop"];
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
	    Prompt: string;
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
	        this.Prompt = source["Prompt"];
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
	export class Config {
	    UI: UIConfig;
	    Backend: string;
	    LocalAI: LocalAIConfig;
	    OpenAI: OpenAIConfig;
	    AnythingLLM: AnythingLLMConfig;
	    Ollama: OllamaConfig;
	    Mistral: MistralConfig;
	    Anthropic: AnthropicConfig;
	    CallOptions: CallOptionsConfig;
	    Printer: PrinterConfig;
	    LogLevel: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UI = this.convertValues(source["UI"], UIConfig);
	        this.Backend = source["Backend"];
	        this.LocalAI = this.convertValues(source["LocalAI"], LocalAIConfig);
	        this.OpenAI = this.convertValues(source["OpenAI"], OpenAIConfig);
	        this.AnythingLLM = this.convertValues(source["AnythingLLM"], AnythingLLMConfig);
	        this.Ollama = this.convertValues(source["Ollama"], OllamaConfig);
	        this.Mistral = this.convertValues(source["Mistral"], MistralConfig);
	        this.Anthropic = this.convertValues(source["Anthropic"], AnthropicConfig);
	        this.CallOptions = this.convertValues(source["CallOptions"], CallOptionsConfig);
	        this.Printer = this.convertValues(source["Printer"], PrinterConfig);
	        this.LogLevel = source["LogLevel"];
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
	
	export class LLMMessage {
	    Role: string;
	    Content: string;
	
	    static createFrom(source: any = {}) {
	        return new LLMMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Role = source["Role"];
	        this.Content = source["Content"];
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

}

