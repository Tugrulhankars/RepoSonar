/* eslint-disable @typescript-eslint/no-unused-vars */
import React, { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../redux/store/store";
import { searchRepos } from "../redux/searchSlice";
import { SearchBar } from "../components/SearchBar";
import { Pagination } from "../components/Pagination";
import { streamAllRepositories } from "../redux/syncSlice";

export const HomePage: React.FC = () => {
  const dispatch = useAppDispatch();
  const { results, totalHits, isLoading, error } = useAppSelector((state) => state.search);

  // Sayfa ilk yüklendiğinde varsayılan popüler repoları getir
  useEffect(() => {
    dispatch(searchRepos({ query: "", page: 1, size: 10 }));
  }, [dispatch]);

  const handleStreamAllRepositories=()=>{
    dispatch(streamAllRepositories())
  }
  return (
    <div
      style={{
        maxWidth: 900,
        margin: "40px auto",
        padding: "0 20px",
        fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif",
      }}
    >
      <header style={{ marginBottom: "25px" }}>
        <h1 style={{ fontSize: "28px", color: "#1f2328", marginBottom: "8px" }}>
          🐙 GitHub Repository Search
        </h1>
        <p style={{ color: "#656d76", fontSize: "14px" }}>
          Elasticsearch destekli hızlı arama, filtreleme ve otomatik tamamlama motoru.
        </p>
      </header>

      {/* Arama ve Filtre Bölümü */}
      <SearchBar />
      <button onClick={handleStreamAllRepositories} style={{
            padding: "10px 24px",
            backgroundColor: "#2da44e",
            color: "#fff",
            border: "none",
            borderRadius: "6px",
            fontWeight: 600,
            cursor: "pointer",
          }}>Stream All Repositories</button>

      {/* Bilgi ve Durum Metni */}
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: "15px" }}>
        <span style={{ fontSize: "14px", color: "#656d76" }}>
          {isLoading
            ? "Sonuçlar aranıyor..."
            : `Toplam ${totalHits.toLocaleString()} sonuç listelendi`}
        </span>
      </div>

      {/* Hata Durumu */}
      {error && (
        <div
          style={{
            padding: "12px",
            backgroundColor: "#ffebe9",
            color: "#cf222e",
            borderRadius: "6px",
            marginBottom: "20px",
          }}
        >
          {error}
        </div>
      )}

      {/* Sonuç Listesi */}
      <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
        {results.map((repo) => (
          <div
            key={repo.id}
            style={{
              border: "1px solid #d0d7de",
              borderRadius: "8px",
              padding: "18px",
              backgroundColor: "#ffffff",
              boxShadow: "0 1px 3px rgba(0,0,0,0.02)",
            }}
          >
            {/* Repo Başlığı & Yıldız */}
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
              <a
                href={repo.html_url}
                target="_blank"
                rel="noreferrer"
                style={{
                  fontSize: "18px",
                  fontWeight: 600,
                  color: "#0969da",
                  textDecoration: "none",
                }}
              >
                {repo.full_name}
              </a>
              <span
                style={{
                  fontSize: "13px",
                  fontWeight: 600,
                  color: "#57606a",
                  border: "1px solid #d0d7de",
                  padding: "3px 8px",
                  borderRadius: "12px",
                  backgroundColor: "#f6f8fa",
                }}
              >
                ⭐ {repo.stargazers_count.toLocaleString()}
              </span>
            </div>

            {/* Açıklama */}
            <p style={{ color: "#57606a", fontSize: "14px", margin: "10px 0 14px 0", lineHeight: 1.45 }}>
              {repo.description || "Açıklama girilmemiş."}
            </p>

            {/* Dil ve Topics Etiketleri */}
            <div style={{ display: "flex", gap: "8px", flexWrap: "wrap", alignItems: "center" }}>
              {repo.language && (
                <span
                  style={{
                    fontSize: "12px",
                    fontWeight: 500,
                    color: "#24292f",
                    backgroundColor: "#f6f8fa",
                    padding: "3px 8px",
                    borderRadius: "12px",
                    border: "1px solid #d0d7de",
                  }}
                >
                  ● {repo.language}
                </span>
              )}
              {repo.topics?.map((topic, index) => (
                <span
                  key={index}
                  style={{
                    fontSize: "12px",
                    color: "#0969da",
                    backgroundColor: "#ddf4ff",
                    padding: "2px 8px",
                    borderRadius: "12px",
                    fontWeight: 500,
                  }}
                >
                  #{topic}
                </span>
              ))}
            </div>
          </div>
        ))}
      </div>

      {/* Sayfalama Kontrolleri */}
      <Pagination pageSize={10} />
    </div>
  );
};