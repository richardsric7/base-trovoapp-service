"use client";
import React, { useState } from "react";
import styled from "styled-components";
import PrimaryButton from "@/components/PrimaryButton";
import { FaPlus } from "react-icons/fa6";

import { SearchBar } from "@/components";
import SuccessMessage from "@/components/SuccessMessage";
import AssetCurationFilter from "../assetcuration/components/AssetCurationFilter";

import AddInitiatorsModal from "./components/AddInitiatorsModal";
import InitiatorsTable from "./components/InitiatorsTable";

const Initiators = () => {
  const [add, setAdd] = useState<boolean>(false);
  const [added, setAdded] = useState<boolean>(false);
  const [username, setUsername] = useState<string>("");
  const [search, setSearch] = useState<string>("");
  const [searchTrigger, setSearchTrigger] = useState<string>("");

  const handleSearchInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearch(e.target.value);
    if (e.target.value.trim() === "") {
      setSearchTrigger("");
    }
  };

  const handleSearch = () => {
    setSearchTrigger(search); // <-- THIS triggers a table refetch
  };

  return (
    <PageContainer>
      <Header>
        <TitleSection>
          <Title> Minting Initiators</Title>
        </TitleSection>

        <PrimaryButton
          buttonStyle={{
            width: "auto",
            padding: "10px 20px",
            display: "flex",
            alignItems: "center",
            gap: "4px",
          }}
          onClick={() => setAdd(!add)}
        >
          <FaPlus /> Initiators
        </PrimaryButton>
      </Header>

      <FiltersSection>
        <SearchBar
          value={search}
          onChange={handleSearchInputChange}
          handleSearch={handleSearch}
        />
        <AssetCurationFilter />
      </FiltersSection>

      <InitiatorsTable search={searchTrigger} />
      {add && (
        <AddInitiatorsModal
          add={add}
          setAdd={setAdd}
          setAdded={setAdded}
          username={username}
          setUsername={setUsername}
        />
      )}

      {added && (
        <SuccessMessage
          isOpen={added}
          setIsOpen={setAdded}
          heading="Success"
          message={`You have successfully added ${username} as an`}
          email="Initiator"
        />
      )}
    </PageContainer>
  );
};

export default Initiators;
const PageContainer = styled.section``;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const TitleSection = styled.div`
  flex-grow: 1;
`;

const Title = styled.h1`
  font-size: 24px;
  font-weight: 700;
  // line-height: 28px;
  margin: 0;
  color: #00225a;
`;

const FiltersSection = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 25px;
  margin: 10px 0 30px 0;
`;
