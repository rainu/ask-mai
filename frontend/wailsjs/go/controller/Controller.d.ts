// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {config} from '../models';
import {controller} from '../models';

export function AppMounted():Promise<void>;

export function GetApplicationConfig():Promise<config.Config>;

export function LLMAsk(arg1:controller.LLMAskArgs):Promise<string>;

export function LLMInterrupt():Promise<void>;

export function LLMWait():Promise<string>;
