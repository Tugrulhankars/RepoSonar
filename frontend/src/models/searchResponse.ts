import type { Repository } from "./repository";

export interface SearchResponse {
    query: string;
    total_hits: number;
    page: number;
    size: number;
    results: Repository[];
}