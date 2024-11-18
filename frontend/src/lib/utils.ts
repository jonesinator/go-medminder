import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

const backendBaseUrl = process.env.NEXT_PUBLIC_BACKEND_URL || "/"

export function apiUrl(path: string): string {
    return backendBaseUrl + path
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}