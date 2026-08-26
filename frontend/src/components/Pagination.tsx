import React from "react";
import { useAppDispatch, useAppSelector } from "../redux/store/store";
import { searchRepos, setPage } from "../redux/searchSlice";

interface PaginationProps {
  pageSize?: number;
}

export const Pagination: React.FC<PaginationProps> = ({ pageSize = 10 }) => {
  const dispatch = useAppDispatch();
  const { page, totalHits, query, language, minStars, isLoading } = useAppSelector(
    (state) => state.search
  );

  const totalPages = Math.ceil(totalHits / pageSize);

  if (totalPages <= 1) return null;

  const handlePageChange = (newPage: number) => {
    if (newPage < 1 || newPage > totalPages || isLoading) return;

    dispatch(setPage(newPage));
    dispatch(
      searchRepos({
        query,
        language,
        stars: minStars,
        page: newPage,
        size: pageSize,
      })
    );

    // Sayfa değiştiğinde sayfanın en üstüne yumuşakça kaydır
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        gap: "12px",
        marginTop: "30px",
        paddingBottom: "40px",
      }}
    >
      <button
        onClick={() => handlePageChange(page - 1)}
        disabled={page === 1 || isLoading}
        style={{
          padding: "8px 16px",
          border: "1px solid #d0d7de",
          borderRadius: "6px",
          backgroundColor: page === 1 ? "#f6f8fa" : "#fff",
          cursor: page === 1 ? "not-allowed" : "pointer",
          fontWeight: 500,
        }}
      >
        &larr; Önceki
      </button>

      <span style={{ fontSize: "14px", color: "#57606a" }}>
        Sayfa <strong>{page}</strong> / <strong>{totalPages}</strong>
      </span>

      <button
        onClick={() => handlePageChange(page + 1)}
        disabled={page === totalPages || isLoading}
        style={{
          padding: "8px 16px",
          border: "1px solid #d0d7de",
          borderRadius: "6px",
          backgroundColor: page === totalPages ? "#f6f8fa" : "#fff",
          cursor: page === totalPages ? "not-allowed" : "pointer",
          fontWeight: 500,
        }}
      >
        Sonraki &rarr;
      </button>
    </div>
  );
};