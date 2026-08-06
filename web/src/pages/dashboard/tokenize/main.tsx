import { useNavigate } from "react-router-dom";
import Header from "../../../components/header";
import { FaCirclePlus } from "react-icons/fa6";

export const Tokenize = () => {
    const navigate = useNavigate();
  return (
    <div>
      <Header />
      <div className="max-w-full p-5">
        <div className="space-y-8">
          <div className="bg-[#004988] text-white rounded-2xl p-6 text-center shadow-md bg-[url('/images/topographicIcon.png')] bg-cover bg-center">
            <h2 className="text-xl font-semibold mb-2">
              Unlock Liquidity From Fixed Assets
            </h2>
            <p className="text-sm mb-4 font-light">
              Digitally fragment any asset of value to unlock capital from
              multiple contributors. Use the power of tokenization to give more
              people access to investment opportunities previously unavailable
              to them. Tokenize land, real estate, infrastructure projects,
              commodities, natural resources and other qualifying assets.
            </p>
            <button className="mt-2 inline-flex items-center gap-2 px-4 py-2 bg-[#FFFFFF33] text-[#FFFFFF] font-semibold rounded-lg hover:bg-gray-100 transition">
              <FaCirclePlus />
              Tokenize Asset
            </button>
          </div>

          <div className="bg-[#EDF3F9] rounded-xl p-8 text-center">
            <div className="flex justify-center mb-6">
              <img
                src="/images/emptyBox.png"
                alt="Tokenize Graphic"
                className="w-397 h-185 object-contain"
              />
            </div>

            <p className="text-sm text-[#0E1F51] max-w-2xl mx-auto mb-4">
              Tokenized assets that are for public offering are securities and
              MUST be approved by the Securities and Exchange Commission (SEC).
            </p>
            <p className="text-sm text-[#0E1F51] max-w-2xl mx-auto mb-4">
              Only assets approved by the SEC can be tokenized and offered to
              the public on the Trovo tokenization platform.
            </p>
            <p className="text-sm text-[#0E1F51] max-w-2xl mx-auto mb-6">
              Asset issued by private entities and are for private offering are
              not under the purview of the SEC.
            </p>

            <p className="text-sm font-medium text-[#0E1F51] mb-4">
              <span className="font-semibold">
                Are you an asset owner and want to tokenize your asset to
                attract capital from around the world?
              </span>
            </p>

            <button onClick={()=> navigate('/dashboard/tokenize/apply')} className="inline-flex items-center gap-2 px-5 py-2 bg-[#004988] text-white rounded-lg hover:bg-[#00285B] transition">
              <FaCirclePlus size={20} />
              Proceed to Tokenize Asset
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
