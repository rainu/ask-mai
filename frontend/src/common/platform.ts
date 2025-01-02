export const IsWindows: boolean = navigator.userAgent.includes("Windows")
export const PathSeparator: string = IsWindows ? '\\' : '/'