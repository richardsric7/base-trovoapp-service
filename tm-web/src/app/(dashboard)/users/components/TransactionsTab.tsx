import React, { useMemo, useState } from "react";
import CustomTable from "@/components/CustomTable";
import styled from "styled-components";
import { FaArrowDown, FaArrowUp, FaMoneyBill } from "react-icons/fa6";
import { FaExchangeAlt } from "react-icons/fa";
import { SearchBar } from "@/components";
import TransactionsFilter from "./TransactionsFilter";
import Pagination from "@/components/CustomPagination";
import TransactionModal from "./TransactionModal";
import { useGetPaymentHistoryQuery } from "@/redux/api/users";
import { useParams } from "next/navigation";
import SelectedFiltersComponent from "@/components/SelectedFilter";
interface FilterState {
  from?: string;
  category?: string;
  date?: string;
  [key: string]: string | undefined;
}

const TransactionsTab = () => {
  const params = useParams();
  const username = Array.isArray(params.username)
    ? params.username[0]
    : params.username;

  const [currentPage, setCurrentPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [openModal, setOpenModal] = useState<boolean>(false);
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [selectedTransaction, setSelectedTransaction] = useState(null);
  const [filters, setFilters] = useState<FilterState>({});

  const { data, isLoading, isFetching, refetch } = useGetPaymentHistoryQuery({
    page: currentPage,
    pageSize,
    search: searchQuery,
    from: username,
    ...Object.fromEntries(
      Object.entries(filters).filter(([_, value]) => value !== undefined)
    ),
  });

  const columns = [
    {
      title: "Transaction",
      dataIndex: "category",
      render: (_: any, record: any) => {
        let categoryIcon;
        let iconBgColor;
        let textColor;

        switch (record.category) {
          case "Sent":
            categoryIcon = <FaArrowUp />;
            iconBgColor = "#BE380033";
            textColor = "#BE3800";
            break;
          case "PAYMENT":
            categoryIcon = <FaArrowDown />;
            iconBgColor = "#00A85933";
            textColor = "#00A859";
            break;
          case "SWAP ":
          case "SWAP ":
          case "SWAP":
            categoryIcon = <FaExchangeAlt />;
            iconBgColor = "#ACD1EF";
            textColor = "#007CDF";
            break;
          case "Tokenization fee":
            categoryIcon = <FaMoneyBill />;
            iconBgColor = "#ACD1EF";
            textColor = "#007CDF";
            break;
          default:
            categoryIcon = <FaArrowDown />;
            iconBgColor = "gray";
            textColor = "gray";
        }

        return (
          <CategoryWrapper>
            <CategoryIcon iconBgColor={iconBgColor} textColor={textColor}>
              {categoryIcon}
            </CategoryIcon>
            <div>
              <CategoryText>{record.category}</CategoryText>
              <DateText>{record.date}</DateText>
            </div>
          </CategoryWrapper>
        );
      },
      key: "category",
    },
    {
      title: "Asset",
      dataIndex: "assetCode",
      key: "assetCode",
      render: (_: any, record: any) => {
        const txType = record.transactionType?.trim() || "";

        // Check for SWAP
        const isSwap = txType.toUpperCase().startsWith("SWAP ");
        const swapPart = isSwap ? txType.substring(5).trim() : ""; // removes "SWAP "

        // Format: XBN>TROV --> TROV/XBN
        const swapDirection =
          isSwap && swapPart.includes(">")
            ? swapPart.split(">").reverse().join("/")
            : "";

        return (
          <WalletName>
            <Avatar />
            {isSwap ? swapDirection || "--" : record.assetCode || "--"}
          </WalletName>
        );
      },
    },
    // {
    //   title: "Asset",
    //   dataIndex: "assetCode",
    //   render: (_: any, record: any) => {
    //     return (
    //       <WalletName>
    //         {" "}
    //         <Avatar />
    //         {record.assetCode || "--"}
    //       </WalletName>
    //     );
    //   },
    //   key: "assetCode",
    // },
    {
      title: "From",
      dataIndex: "from",
      render: (_: any, record: any) => {
        if (!record.from) return <WalletName>--</WalletName>;
        const match = record.from?.match(/\[(.*?)\]/);
        const nameOnly = match ? match[1] : record.from;
        return (
          <WalletName>
            <Avatar />
            {nameOnly}
          </WalletName>
        );
      },
    },

    {
      title: "To",
      dataIndex: "to",
      render: (_: any, record: any) => {
        if (!record.to) return <WalletName>--</WalletName>;
        const match = record.to?.match(/\[(.*?)\]/);
        const nameOnly = match ? match[1] : record.to;
        return (
          <WalletName>
            <Avatar />
            {nameOnly}
          </WalletName>
        );
      },
    },
    {
      title: " Amount",
      dataIndex: "amount",
      key: "amount",
      render: (_: any, record: any) => {
        if (!record.amount) {
          return <Amount color="#BE3800">--</Amount>; // fallback for missing
        }

        const amountValue = parseFloat(record.amount);
        const category = record.category?.toUpperCase() || "";

        let color = ""; // Default red
        if (category === "PAYMENT") {
          color = "#00A859"; // Green
        } else if (category.startsWith("SWAP")) {
          color = "#007CDF"; // Blue for swap
        }

        return (
          <Amount color={color}>
            {amountValue.toFixed(2)} <span>{record.assetCode}</span>
          </Amount>
        );
      },
    },

    {
      title: "",
      dataIndex: "view",
      key: "view",
      render: (_: any, record: any) => {
        return (
          <DetailsLink
            onClick={() => {
              setOpenModal(!openModal), setSelectedTransaction(record);
            }}
          >
            View
          </DetailsLink>
        );
      },
    },
  ];
  const formatDate = (date: string) =>
    new Date(date).toLocaleDateString("en-GB", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    });

  const formatAmount = (amount: number | string) => {
    const num = typeof amount === "string" ? parseFloat(amount) : amount;
    if (isNaN(num)) return "--";
    return new Intl.NumberFormat("en-US", {
      minimumFractionDigits: 2, // Always show at least 2
      maximumFractionDigits: 7, // Show up to 7
    }).format(num);
  };

  const dataSource = useMemo(() => {
    if (isLoading || !Array.isArray(data?.data?.data)) {
      return [];
    }
    return data.data.data.map((record: any) => ({
      key: record.transactionId,

      amount: formatAmount(record.amount),

      description: record.memo || "No description",
      category: record.transactionType?.split(" ")[0] || record.transactionType,

      date: formatDate(record.transactionDate),

      transactionId: record.transactionId,
      to: record.to,
      from: record.from,
      assetCode: record.assetCode,
      fromPublicKey: record.fromPublicKey,
      toPublicKey: record.toPublicKey,
      transactionType: record.transactionType,
    }));
  }, [data, isLoading]);

  const handlePageChange = (pageNumber: number) => {
    setCurrentPage(pageNumber);
  };
  const handleSearch = () => {
    setSearchQuery(searchQuery);
    setCurrentPage(1);
    refetch();
    console.log("Search Query:", searchQuery);
  };

  if (isLoading) {
    return (
      <Container>
        <Text>Loading Transactions...</Text>
      </Container>
    );
  }

  const handleFilters = (newFilters: FilterState) => {
    setFilters(newFilters);
    setCurrentPage(1);
    refetch();
  };

  const handleRemoveFilter = (key: keyof FilterState) => {
    const newFilters = { ...filters };
    delete newFilters[key];
    setFilters(newFilters);
    refetch();
  };

  const transformedFilters = Object.fromEntries(
    Object.entries(filters).filter(([_, value]) => value !== undefined)
  ) as Record<string, string>;

  return (
    <Container>
      {dataSource.length === 0 ? (
        <Text> No Transactions yet! </Text>
      ) : (
        <>
          <Content>
            <div>
              <SearchBar
                customWidth="320px"
                placeholder="Start search..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                handleSearch={handleSearch}
              />
            </div>
            <TransactionsFilter
              onApplyFilters={handleFilters}
              currentFilters={filters}
            />
          </Content>
          <SelectedFiltersComponent
            selectedFilters={transformedFilters}
            handleRemoveFilter={handleRemoveFilter}
          />
          <CustomTable columns={columns} dataSource={dataSource} />
          <PaginationWrapper>
            <Pagination
              currentPage={currentPage}
              totalCount={data?.data?.total || 0}
              pageSize={pageSize}
              onPageSizeChange={setPageSize}
              onPageChange={handlePageChange}
              isFetching={isFetching}
            />
          </PaginationWrapper>
          {openModal && (
            <TransactionModal
              openModal={openModal}
              setOpenModal={setOpenModal}
              transaction={selectedTransaction}
            />
          )}
        </>
      )}
    </Container>
  );
};

export default TransactionsTab;

const Container = styled.section`
  background-color: #fff;
  padding: 20px;
  margin-bottom: 20px;
  border-radius: 24px 24px 0 0;
  table {
    width: 100%;
    table-layout: fixed;
  }
  th:nth-child(1),
  td:nth-child(1) {
    width: 12%;
  }

  th:nth-child(2),
  td:nth-child(2) {
    width: 14%;
  }

  th:nth-child(3),
  td:nth-child(3) {
    width: 16%;
  }
  th:nth-child(4),
  td:nth-child(4) {
    width: 14%;
  }

  th:nth-child(5),
  td:nth-child(5) {
    width: 10%;
  }

  th:nth-child(6),
  td:nth-child(6) {
    width: 6%;
  }
`;

const Content = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
`;

const WalletName = styled.p`
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 2px;
`;

const CategoryWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;

const CategoryIcon = styled.div<{ iconBgColor: string; textColor: string }>`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: ${(props) => props.iconBgColor};
  color: ${(props) => props.textColor};
`;

const CategoryText = styled.span`
  color: #00225a;
  font-weight: 600;
`;

const Amount = styled.span<{ color: string }>`
  color: ${(props) => props.color};
  font-size: 12px;
  font-weight: 500;
  line-height: 16px;
`;

const PaginationWrapper = styled.div`
  margin-top: 20px;
`;

const Text = styled.p`
  font-size: 16px;
  font-weight: 500;
  line-height: 16px;
  color: #007cdf;
  text-align: center;
  margin: 20px auto;
`;

const DetailsLink = styled.p`
  color: #007cdf;
  text-decoration: underline;
`;
const Avatar = styled.div`
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: rebeccapurple;
  max-width: 36px;
`;

const DateText = styled.p``;
