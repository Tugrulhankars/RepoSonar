/* eslint-disable @typescript-eslint/no-explicit-any */
import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import axios from "axios";




const axiosClient=axios.create({
    baseURL:"http://localhost:8082/api/v1",
    timeout:10000,
    headers:{
        "Content-Type":"application/json"
    }
})




export const streamAllRepositories= createAsyncThunk<any, void, { rejectValue: string }>(
    "sync/streamAllRepositories",
    async (_, { rejectWithValue }) => {
      try {
        const response = await axiosClient.post("/streamAllRepositories");
        return response.data;
      } catch (err) {
        if (axios.isAxiosError(err) && err.response) {
          return rejectWithValue(err.response.data?.error || "Arama sırasında bir hata oluştu");
        }
        return rejectWithValue("Sunucuya ulaşılamadı");
      }
    }
  );


export const syncSlice=createSlice({
    name:"sync",
    initialState:{},
    reducers:{},
    extraReducers:(builder)=>{
        builder.addCase(streamAllRepositories.fulfilled,(state,action)=>{
            return action.payload
        })
    }
})


export default syncSlice.reducer