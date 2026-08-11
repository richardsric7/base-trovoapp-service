import Image from "next/image";
import image from "../../../../public/img/new-website/problemwearesolving.png";
import { Lock, Lock1, ReceiptMinus, TrendDown } from "iconsax-reactjs";

const ProblemWeAreSolving = () => {
  const dataContent = [
    {
      heading: "Projects remain underfunded",
      text: "Capital cannot reach the assets that need it most",
      icon: <ReceiptMinus size="32" variant="Bold" color=" #007CDF" />,
    },
    {
      heading: "Investment access remains limited",
      text: "Opportunities locked behind high minimums and geography",
      icon: <Lock size="32" variant="Bold" color="#007CDF" />,
    },
    {
      heading: "Capital allocation becomes inefficient",
      text: "Wealth concentrates while productive assets go underserved",
      icon: <TrendDown size="32" variant="Bold" color="#007CDF" />,
    },
  ];
  return (
    <section className="py-20 md:py-[100px] bg-white flex flex-col md:flex-row gap-6 lg:gap-10 xl:gap-16 justify-between px-4 md:px-12 lg:px-[90px] items-stretch">
      <div className="flex flex-col gap-6 w-full md:w-[48%] lg:w-[45%] xl:w-[40%]">
        <div className="">
          <p className="text-primary font-bold tracking-[0.1em] uppercase mb-4 font-matahari text-sm sm:text-base">
            THE PROBLEM WE ARE SOLVING
          </p>
          <h2 className="text-4xl sm:text-5xl font-extrabold leading-[1.1] mb-6 text-secondary font-matahari">
            The Capital <br /> Coordination Gap
          </h2>
          <p className="text-gray-600 leading-[160%] text-sm sm:text-base font-normal">
            Across emerging markets, productive assets often struggle to access
            structured capital despite strong investor demand. Traditional
            capital formation systems remain fragmented, inaccessible, and
            inefficient for many asset owners and investors. As a result:
          </p>
        </div>

        <div className="flex flex-col gap-4">
          {dataContent.map((item, index) => (
            <div
              className="flex gap-2 bg-light3 px-[16px] pt-[16px] pb-[8px] rounded-md"
              key={index}
            >
              <div className="bg-white rounded-full w-[48px] h-[48px] flex items-center justify-center flex-shrink-0">
                {item.icon}
              </div>
              <div>
                <h2 className="text-lg font-semibold mb-2">{item.heading}</h2>
                <p className="text-trovo-black font-normal ">{item.text}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="bg-[#F2F6F9] rounded-md p-4 w-full md:w-[48%] lg:w-[50%] xl:w-[50%] flex flex-col justify-center">
        <Image
          src={image}
          alt="trovotech-problem-we-are-solving"
          width={539}
          height={640}
          className="w-full h-auto rounded-md"
          priority
        />
      </div>
    </section>
  );
};

export default ProblemWeAreSolving;
