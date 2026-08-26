/* eslint-disable react-hooks/set-state-in-effect */
import React, { useEffect, useState, useRef } from "react";
import { useAppDispatch, useAppSelector } from "../redux/store/store";
import {
  setQuery,
  setLanguage,
  setMinStars,
  getSuggestions,
  clearSuggestions,
  searchRepos,
  setPage,
} from "../redux/searchSlice";
import { useDebounce } from "../hooks/useDebounce";

export const SearchBar: React.FC = () => {
  const dispatch = useAppDispatch();
  const { query, language, minStars } = useAppSelector((state) => state.search);
  const { suggestions } = useAppSelector((state) => state.search);

  const [showDropdown, setShowDropdown] = useState(false);
  const debouncedQuery = useDebounce(query, 300);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Debounced input değiştikçe autocomplete çek
  useEffect(() => {
    if (debouncedQuery.trim().length >= 2) {
      dispatch(getSuggestions(debouncedQuery));
      setShowDropdown(true);
    } else {
      dispatch(clearSuggestions());
      setShowDropdown(false);
    }
  }, [debouncedQuery, dispatch]);

  // Dropdown dışına tıklandığında menüyü kapat
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setShowDropdown(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleSearch = (searchQuery: string = query) => {
    setShowDropdown(false);
    dispatch(setPage(1)); // Yeni aramada sayfayı 1'e sıfırla
    dispatch(
      searchRepos({
        query: searchQuery,
        language,
        stars: minStars,
        page: 1,
        size: 10,
      })
    );
  };

  const handleSelectSuggestion = (suggestedText: string) => {
    dispatch(setQuery(suggestedText));
    handleSearch(suggestedText);
  };

  return (
    <div style={{ position: "relative", marginBottom: "25px" }} ref={dropdownRef}>
      <div style={{ display: "flex", gap: "10px", flexWrap: "wrap" }}>
        {/* Arama Inputu */}
        <input
          type="text"
          placeholder="Repo adı, konu (#topic) veya açıklama ara..."
          value={query}
          onChange={(e) => dispatch(setQuery(e.target.value))}
          onKeyDown={(e) => e.key === "Enter" && handleSearch()}
          style={{
            flex: 2,
            minWidth: "220px",
            padding: "12px",
            fontSize: "15px",
            border: "1px solid #d1d5db",
            borderRadius: "6px",
            outline: "none",
          }}
        />

       <select
  value={language}
  onChange={(e) => {
    const selectedLang = e.target.value;
    dispatch(setLanguage(selectedLang));
    dispatch(setPage(1));

    dispatch(
      searchRepos({
        query: query,
        language: selectedLang,
        stars: minStars,
        page: 1,
        size: 10,
      })
    );
  }}
  style={{
    padding: "10px 14px",
    borderRadius: "6px",
    border: "1px solid #d1d5db",
    backgroundColor: "#fff",
    cursor: "pointer",
  }}
>
  <option value="">Tüm Diller</option>
  <option value="Go">Go</option>
  <option value="Java">Java</option>
  <option value="Python">Python</option>
  <option value="Rust">Rust</option>
</select>

        {/* Min Yıldız Filtresi */}
        <select
  value={minStars}
  onChange={(e) => {
    const selectedStars = Number(e.target.value);
    
    // 1. Redux state'ini güncelle
    dispatch(setMinStars(selectedStars));
    dispatch(setPage(1));

    // 2. Anında yeni yıldız filtresiyle arama isteği at
    dispatch(
      searchRepos({
        query: query,
        language: language,
        stars: selectedStars, // Seçilen anlık değer gönderilir
        page: 1,
        size: 10,
      })
    );
  }}
  style={{
    padding: "10px 14px",
    borderRadius: "6px",
    border: "1px solid #d1d5db",
    backgroundColor: "#fff",
    cursor: "pointer",
  }}
>
  <option value={0}>Tüm Yıldızlar</option>
  <option value={500}>⭐ 500+</option>
  <option value={5000}>⭐ 5,000+</option>
  <option value={20000}>⭐ 20,000+</option>
  <option value={50000}>⭐ 50,000+</option>
  <option value={100000}>⭐ 100,000+</option>
  <option value={200000}>⭐ 200,000+</option>
  <option value={500000}>⭐ 500,000+</option>
  <option value={1000000}>⭐ 1,000,000+</option>
</select>

        {/* Ara Butonu */}
        <button
          onClick={() => handleSearch()}
          style={{
            padding: "10px 24px",
            backgroundColor: "#2da44e",
            color: "#fff",
            border: "none",
            borderRadius: "6px",
            fontWeight: 600,
            cursor: "pointer",
          }}
        >
          Ara
        </button>
      </div>

      {/* Autocomplete Dropdown Listesi */}
      {showDropdown && suggestions.length > 0 && (
        <ul
          style={{
            position: "absolute",
            top: "48px",
            left: 0,
            width: "50%",
            minWidth: "250px",
            backgroundColor: "#fff",
            border: "1px solid #e1e4e8",
            borderRadius: "6px",
            listStyle: "none",
            margin: 0,
            padding: "6px 0",
            zIndex: 100,
            boxShadow: "0 8px 24px rgba(149,157,165,0.2)",
          }}
        >
          {suggestions.map((item, idx) => (
            <li
              key={idx}
              onClick={() => handleSelectSuggestion(item)}
              style={{
                padding: "10px 16px",
                cursor: "pointer",
                fontSize: "14px",
                borderBottom: idx !== suggestions.length - 1 ? "1px solid #f6f8fa" : "none",
              }}
              onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = "#f3f4f6")}
              onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = "#fff")}
            >
              🔎 {item}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};