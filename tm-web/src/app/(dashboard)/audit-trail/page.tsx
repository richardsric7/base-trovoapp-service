"use client";
import { SearchBar } from "@/components";
import { IAuditTrailParams } from "@/redux/api/auditTrail/interface";
import { useState } from "react";
import styled from "styled-components";
import AuditFilter from "./components/AuditFilter";
import AuditTable from "./components/AuditTable";

const AuditTrailPage = () => {
  // The committed filters that drive the query.
  const [filters, setFilters] = useState<IAuditTrailParams>({});
  // The uncommitted search box value (applied on submit).
  const [searchTerm, setSearchTerm] = useState("");

  const applySearch = () => {
    setFilters((prev) => ({ ...prev, search: searchTerm || undefined }));
  };

  const applyFilters = (next: { status?: string; date?: string }) => {
    setFilters((prev) => ({
      ...prev,
      status: next.status || undefined,
      date: next.date || undefined,
    }));
  };

  return (
    <Container>
      <Title>Audit Trail</Title>
      <SubHeader>
        <SearchBar
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          handleSearch={applySearch}
        />
        <AuditFilter onApply={applyFilters} />
      </SubHeader>
      <AuditTable filters={filters} />
    </Container>
  );
};

export default AuditTrailPage;

const Container = styled.section`
  background: #ffffff;
  padding: 32px;

  border-radius: 24px;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  color: #00225a;
`;

const SubHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 18px;
`;
