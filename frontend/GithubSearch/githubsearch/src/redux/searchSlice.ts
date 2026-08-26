import type { SearchFilterParams } from '../models/searchFilterParams';
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unused-vars */
import { createAsyncThunk, createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { Repository } from "../models/repository";
import axios from "axios";
import type { SuggestResponse } from "../models/suggestResponse";
import type { SearchResponse } from "../models/searchResponse";

//const BASE_URL="http://localhost:8082/api/v1";
const axiosClient = axios.create({
  baseURL: "http://localhost:8082/api/v1",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});
interface SearchState {
    // Arama filtreleri
  query: string;
  language: string;
  minStars: number;
  page: number;

  // Autocomplete state
  suggestions: string[];
  isSuggestLoading: boolean;

  // Arama sonuçları state
  results: Repository[];
  totalHits: number;
  isLoading: boolean;
  error: string | null;
}

const initialState:SearchState={
    query: "",
    language: "",
    minStars: 0,
    page: 1,
    suggestions: [],
    isSuggestLoading: false,
    results: [],
    totalHits: 0,
    isLoading: false,
    error: null
}

// 1. Öneri Çekme (Autocomplete) Thunk'ı - Axios ile
export const getSuggestions = createAsyncThunk<string[], string, { rejectValue: string }>(
  "search/getSuggestions",
  async (prefix: string, { rejectWithValue }) => {
    try {
      if (!prefix.trim()) return [];
      
      const response = await axiosClient.get<SuggestResponse>("/suggest", {
        params: { prefix },
      });

      return response.data.suggestions || [];
    } catch (err) {
      if (axios.isAxiosError(err) && err.response) {
        return rejectWithValue(err.response.data?.error || "Öneriler alınamadı");
      }
      return rejectWithValue("Bağlantı hatası oluştu");
    }
  }
);


export const searchRepos=createAsyncThunk<SearchResponse,SearchFilterParams,{rejectValue:string}>(

    "search/searchRepos",
    async (params:SearchFilterParams,{rejectWithValue})=>{
        try {
      const response = await axiosClient.get<SearchResponse>("/search", {
        params: {
          q: params.query || "",
          lang: params.language || "",
          stars: params.stars || 0,
          page: params.page || 1,
          size: params.size || 10,
        },
      });

      return response.data;
    } catch (err) {
      if (axios.isAxiosError(err) && err.response) {
        return rejectWithValue(err.response.data?.error || "Arama sırasında bir hata oluştu");
      }
      return rejectWithValue("Sunucuya ulaşılamadı");
    }
    }
)


export const searchSlice=createSlice({
    name: "search",
    initialState,
    reducers: {
        setQuery: (state, action: PayloadAction<string>) => {
      state.query = action.payload;
    },
    setLanguage: (state, action: PayloadAction<string>) => {
      state.language = action.payload;
    },
    setMinStars: (state, action: PayloadAction<number>) => {
      state.minStars = action.payload;
    },
    setPage: (state, action: PayloadAction<number>) => {
      state.page = action.payload;
    },
    clearSuggestions: (state) => {
      state.suggestions = [];
    },
    },
    extraReducers: (builder)=>{
        // Autocomplete
    builder
      .addCase(getSuggestions.pending, (state) => {
        state.isSuggestLoading = true;
      })
      .addCase(getSuggestions.fulfilled, (state, action) => {
        state.isSuggestLoading = false;
        state.suggestions = action.payload;
      })
      .addCase(getSuggestions.rejected, (state) => {
        state.isSuggestLoading = false;
        state.suggestions = [];
      });

    // Search
    builder
      .addCase(searchRepos.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(searchRepos.fulfilled, (state, action: PayloadAction<SearchResponse>) => {
        state.isLoading = false;
       // Backend'den 'data' veya 'results' olarak gelse bile yakalar:
        state.results = action.payload.data || action.payload.results || [];
        state.totalHits = action.payload.total_hits || 0;
      })
      .addCase(searchRepos.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload || "Bilinmeyen bir hata oluştu";
      });
    }
})


export const { setQuery, setLanguage, setMinStars, setPage, clearSuggestions } = searchSlice.actions;
export default searchSlice.reducer;