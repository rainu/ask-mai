export namespace config {
	
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
	export class OpenAIConfig {
	    APIKey: string;
	    SystemPrompt: string;
	
	    static createFrom(source: any = {}) {
	        return new OpenAIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APIKey = source["APIKey"];
	        this.SystemPrompt = source["SystemPrompt"];
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
	    BackgroundColor: WindowBackgroundColor;
	    StartState: number;
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
	        this.BackgroundColor = this.convertValues(source["BackgroundColor"], WindowBackgroundColor);
	        this.StartState = source["StartState"];
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
	    OpenAI: OpenAIConfig;
	    AnythingLLM: AnythingLLMConfig;
	    Printer: PrinterConfig;
	    LogLevel: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UI = this.convertValues(source["UI"], UIConfig);
	        this.Backend = source["Backend"];
	        this.OpenAI = this.convertValues(source["OpenAI"], OpenAIConfig);
	        this.AnythingLLM = this.convertValues(source["AnythingLLM"], AnythingLLMConfig);
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

