

const ALL_VALUES = [
  {
    number: "01",
    title: "Trust & Governance",
    desc: "We believe sustainable financial systems are built on transparency, accountability, and institutional integrity.",
    row: "top",
  },
  {
    number: "02",
    title: "Productive Capital",
    desc: "We focus on enabling capital to flow toward assets and opportunities that generate real economic value.",
    row: "top",
  },
  {
    number: "03",
    title: "Regulatory Alignment",
    desc: "We build within compliant frameworks designed to support long-term market confidence and stability.",
    row: "top",
  },
  {
    number: "04",
    title: "Inclusive Access",
    desc: "We believe broader participation in investment opportunities can contribute to stronger economic growth and wealth creation.",
    row: "bottom",
  },
  {
    number: "05",
    title: "Long-Term Infrastructure Thinking",
    desc: "We are building systems designed for durability, scalability, and institutional adoption across emerging markets.",
    row: "bottom",
  },
];

const CardItem = ({ value }: { value: (typeof ALL_VALUES)[0] }) => (
  <div className="flex flex-col p-[16px] rounded-[12px] bg-[#F2F6F9] hover:-translate-y-1.5 transition-all duration-300 group h-full">
    <div className="text-primary text-[20px] font-bold mb-1 group-hover:scale-105 transition-transform duration-300">
      {value.number}
    </div>
    <h3 className="text-xl sm:text-2xl font-extrabold text-secondary font-matahari mb-6 leading-tight group-hover:text-primary transition-colors duration-300">
      {value.title}
    </h3>
    <p className="text-gray-600 text-sm leading-[170%] font-sans font-normal">
      {value.desc}
    </p>
  </div>
);

const Values = () => {
  const topRow = ALL_VALUES.filter((v) => v.row === "top");
  const bottomRow = ALL_VALUES.filter((v) => v.row === "bottom");

  return (
    <section className="py-20 md:py-[120px] bg-white">
      <div className=" mx-auto px-5 sm:px-8 md:px-12 lg:px-20">

        {/* Header */}
        <div className="text-center max-w-[750px] mx-auto mb-16 flex flex-col items-center gap-[20px]">
          <p className="text-primary font-bold tracking-[0.15em] uppercase text-xs sm:text-sm font-matahari">
            OUR CORE VALUES
          </p>
          <h2 className="text-3xl sm:text-4xl md:text-5xl lg:text-[48px] font-extrabold leading-[1.14]  text-secondary font-matahari tracking-tight">
            What we stand for
          </h2>
          <p className="text-gray-600 text-base sm:text-lg leading-[175%] font-sans font-normal max-w-[650px]">
            These are not aspirations. They are the principles that govern how we build, who we partner with, and what we refuse to compromise on.
          </p>
        </div>

        {/* Values Grid — 6-col base: top row = 3 × col-span-2, bottom row = 2 × col-span-3 */}
        <div className="grid grid-cols-1 md:grid-cols-6 gap-6 lg:gap-2">

          {/* Top row — 3 equal-width rectangular cards */}
          {topRow.map((value, idx) => (
            <div key={idx} className="md:col-span-2">
              <CardItem value={value} />
            </div>
          ))}

          {/* Bottom row — 2 wider rectangular cards */}
          {bottomRow.map((value, idx) => (
            <div key={idx} className="md:col-span-3">
              <CardItem value={value} />
            </div>
          ))}

        </div>

      </div>
    </section>
  );
};

export default Values;

