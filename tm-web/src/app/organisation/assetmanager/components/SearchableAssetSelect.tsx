"use client";

import { Select } from "antd";
import styled from "styled-components";

export interface SearchableAssetOption {
  id: string;
  code: string;
  name: string;
}

interface SearchableAssetSelectProps {
  assets: SearchableAssetOption[];
  value: string;
  onChange: (assetId: string) => void;
  disabled?: boolean;
  loading?: boolean;
  id?: string;
}

const SearchableAssetSelect = ({
  assets,
  value,
  onChange,
  disabled = false,
  loading = false,
  id,
}: SearchableAssetSelectProps) => (
  <StyledSelect
    id={id}
    value={value || undefined}
    onChange={(assetId) => onChange(String(assetId))}
    disabled={disabled}
    loading={loading}
    showSearch
    allowClear
    onClear={() => onChange("")}
    placeholder={loading ? "Loading assets..." : "Search for an asset"}
    notFoundContent={loading ? "Loading assets..." : "No assets found"}
    options={assets.map((asset) => ({
      value: asset.id,
      label: asset.name ? `${asset.code} — ${asset.name}` : asset.code,
      searchText: `${asset.code} ${asset.name}`.toLowerCase(),
    }))}
    filterOption={(input, option) =>
      String(option?.searchText ?? "").includes(input.trim().toLowerCase())
    }
  />
);

export default SearchableAssetSelect;

const StyledSelect = styled(Select<string>)`
  width: 100%;
  margin-bottom: 6px;

  .ant-select-selector {
    min-height: 46px !important;
    padding: 7px 12px !important;
    border-color: #e0e0e0 !important;
    border-radius: 8px !important;
    box-shadow: none !important;
  }

  .ant-select-selection-search-input {
    height: 44px !important;
  }

  .ant-select-selection-item,
  .ant-select-selection-placeholder {
    line-height: 30px !important;
  }

  &.ant-select-focused .ant-select-selector {
    border-color: #007cdf !important;
  }
`;
