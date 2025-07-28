import { CustomTable } from "../../../components";
import Dropdown from "../../../components/dropdown";
import styles from "./main.module.css";
import { transactionTypeConfig } from "./data";
import Header from "../../../components/header";
import Button from "../../../components/button";
export const History = () => {
  const column = [
    {
      header: "Price",
      key: "price",
      render: (record: any) => {
        console.log(record);

        return (
          <div
            className={`${
              record.type === "Sent" ? "text-red-700" : "text-green-700"
            }`}
          >
            {record.price}
          </div>
        );
      },
    },
    {
      header: "Transaction Type",
      key: "type",
      render: (record: any) => {
        console.log(record);

        return <TransactionType tx={record} />;
      },
    },
    {
      header: "Description",
      key: "description",
    },
    {
      header: "Date",
      key: "date",
    },
  ];
  const dataSource = [
    {
      price: "$100",
      type: "Received",
      description: "Deposit to wallet",
      date: "2023-10-01",
    },
    {
      price: "$50",
      type: "Sent",
      description: "Withdrawal from wallet",
      date: "2023-10-02",
    },
    {
      price: "$50",
      type: "Sent",
      description: "Withdrawal from wallet",
      date: "2023-10-02",
    },
    {
      price: "$50",
      type: "Swap",
      description: "Withdrawal from wallet",
      date: "2023-10-02",
    },
  ];
  return (
    <div>
      <Header />
      <div className="max-w-full p-5">
        <div>
          <h3 className={`font-montserratSemiBold ${styles.headerText}`}>
            Transaction History
          </h3>
          <div className="flex space-x-20 mb-5 w-full">
            <div className="flex space-x-6 flex-1">
              <Dropdown
                label="All wallets"
                options={[{ text: "1", value: "2" }]}
                onSelect={() => {}}
              />
              <Dropdown
                label="All assets"
                options={[{ text: "1", value: "2" }]}
                onSelect={() => {}}
              />
              <Dropdown
                label="All transactions"
                options={[{ text: "1", value: "2" }]}
                onSelect={() => {}}
              />

              <input placeholder="Past Week" onChange={() => {}} type="text" />
            </div>
            <div className="w-1/6">
              <Button label="More Filters" onclick={() => {}} leftIcon={<img src="/images/filterWhite.png" />} />
            </div>
          </div>
          <CustomTable columns={column} data={dataSource} />
        </div>
      </div>
    </div>
  );
};

const TransactionType = ({ tx }: { tx: { type: string } }) => {
  const { icon, color } =
    transactionTypeConfig[tx.type as keyof typeof transactionTypeConfig];
  return (
    <div className="py-4 px-6 flex items-center gap-2 text-sm font-medium text-gray-800">
      <img className="w-7" src={icon} alt="direction" />{" "}
      <span className={`${color}`}>{tx.type}</span>
    </div>
  );
};
