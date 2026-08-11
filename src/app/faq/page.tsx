import PageHeader from "@/component/PageHeader";
import type { ReactNode } from "react";

type Faq = {
  question: string;
  answer: ReactNode;
  open?: boolean;
};

const faqs: Faq[] = [
  {
    question: "What is Trovo App?",
    open: true,
    answer: (
      <p>
        Trovo App is a non-custodial wallet like no other, a multipurpose,
        multi-asset, multi-access wallet built by Trovotech to address practical
        problems faced by existing blockchain wallet users.
      </p>
    ),
  },
  {
    question: "What are the main features of Trovo App?",
    answer: (
      <ul className="list-disc list-outside space-y-2 pl-8">
        <li>Users can recover their wallets even when they lose their secret keys.</li>
        <li>Users can have multiple wallets.</li>
        <li>Users can send, receive, swap and trade multiple assets.</li>
        <li>Wallets can be shared with different access permissions.</li>
        <li>Businesses can maintain their corporate approval processes.</li>
        <li>Store owners can safely receive blockchain payments.</li>
        <li>Real-world assets can be tokenized.</li>
        <li>Voting can be conducted.</li>
      </ul>
    ),
  },
  {
    question: "Can I import an existing wallet from another Bantu wallet?",
    answer: <p>Yes.</p>,
  },
  {
    question: "Why do I need to have a sub-wallet?",
    answer: (
      <>
        <p>
          Sub-wallets are additional wallets users can create under an account
          on the Trovo App. The default wallet is the primary wallet.
        </p>

        <p>
          Sub-wallets allow users to carry out different transactions using one
          Trovo App account instead of creating multiple accounts.
        </p>
      </>
    ),
  },
  {
    question: "What is a Standard Sub Wallet?",
    answer: (
      <>
        <p>
          Standard sub-wallets can be used to hold assets, perform regular
          transactions and share access with other users.
        </p>

        <p>
          Businesses can give staff view-only access to receive payments without
          allowing them to send funds.
        </p>

        <p>
          Authorized users can initiate and approve transactions while other
          stakeholders can view balances and transaction history.
        </p>
      </>
    ),
  },
  {
    question: "What is a Minting/Asset Tokenization Sub Wallet?",
    answer: (
      <>
        <p>
          These wallets are specifically used for issuing or minting new tokens
          on the Bantu Blockchain.
        </p>

        <p>
          The tokens can be ordinary tokens or tokens representing physical and
          intangible assets.
        </p>
      </>
    ),
  },
  {
    question: "What is a Market Making/Trade Sub Wallet?",
    answer: (
      <>
        <p>
          This wallet is used to make markets or carry out trades on the Bantu
          Decentralized Exchange.
        </p>

        <p>
          Users will be able to trade directly on the Bantu DEX through the
          Trovo App.
        </p>

        <p>Only one Market Making/Trade wallet can be created per account.</p>
      </>
    ),
  },
  {
    question: "What is a Bulk Payment Sub Wallet?",
    answer: (
      <p>
        This wallet allows asset tokenizers to make periodic payments to all
        qualified token holders, similar to how dividends are paid in the stock
        market.
      </p>
    ),
  },
  {
    question: "Can I send a transaction to another Bantu wallet using the username?",
    answer: <p>No.</p>,
  },
  {
    question: "Can I upgrade from a gold patron to a diamond patron?",
    answer: <p>Yes.</p>,
  },
  {
    question: "What is shared access and how does it work?",
    answer: (
      <p>
        Shared access enables a Trovo App user to give other users different
        access rights to their wallets.
      </p>
    ),
  },
];

const FaqItem = ({
  question,
  answer,
  open = false,
}: Faq) => {
  return (
    <details className="group" open={open}>
      <summary className="list-none flex cursor-pointer items-center justify-between rounded-[5px] border-l-[3px] border-[#ccc] bg-[#f7f7f7] py-3 pl-4 pr-5 font-semibold text-[#999] transition-all hover:bg-[#f5f5f5] group-open:border-[#007cdf] group-open:text-[#007cdf]">
        <span>{question}</span>

        <span className="block h-2 w-2 origin-[35%] -rotate-45 border-r border-t border-current transition-transform duration-300 group-open:rotate-[135deg]" />
      </summary>

      <div className="space-y-3 pb-2 pt-3 text-[0.9em] leading-relaxed text-[#777]">
        {answer}
      </div>
    </details>
  );
};

const FaqPage = () => {
  return (
    <>
      <PageHeader text="FAQ" />

      <main className="flex-grow">
        <div className="container mx-auto px-4 py-12 md:px-12 lg:px-[90px]">
          <div className="pb-12 text-center">
            <h1 className="text-4xl font-bold tracking-tight text-gray-900 md:text-5xl">
              Frequently Asked{" "}
              <span className="rounded bg-primary px-3 py-1 text-white">
                Questions
              </span>
            </h1>
          </div>

          <section className="mx-auto max-w-4xl">
            <h2 className="mb-8 text-2xl font-normal text-gray-800">
              FAQs for Trovo App
            </h2>

            <div className="space-y-2.5">
              {faqs.map((faq) => (
                <FaqItem
                  key={faq.question}
                  question={faq.question}
                  answer={faq.answer}
                  open={faq.open}
                />
              ))}
            </div>
          </section>
        </div>
      </main>
    </>
  );
};

export default FaqPage;