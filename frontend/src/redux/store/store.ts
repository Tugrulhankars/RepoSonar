import { configureStore } from "@reduxjs/toolkit";
import searchReducer from "../searchSlice";
import { useDispatch, useSelector, type TypedUseSelectorHook } from "react-redux";

export const store=configureStore({
    reducer: {
        search: searchReducer
    }
})



// TypeScript tipleri
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

// Tipli hook'lar (Proje genelinde useDispatch ve useSelector yerine bunları kullanacağız)
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;