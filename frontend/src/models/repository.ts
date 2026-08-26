import type { Owner } from "./owner";

export interface Repository {
    id: number;
    name: string;
    full_name:string;
    description: string;
    html_url: string;
    language: string;
    topics: string[];
    stargazers_count: number;
    forks_count: number;
    owner: Owner;
}