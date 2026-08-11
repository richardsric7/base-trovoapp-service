const STEPS_DATA = [
  {
    number: 1,
    title: "Asset Structuring",
    description:
      "The asset sponsor submits an asset. Legal and financial teams structure it into a compliant financial instrument via a trust or SPV.",
  },
  {
    number: 2,
    title: "Validation and Admission",
    description:
      "Assets undergo legal validation of ownership, financial assessment, risk identification, and regulatory admission review.",
  },
  {
    number: 3,
    title: "Digital Issuance",
    description:
      "Compliant digital tokens are issued representing defined financial interests — income rights, ownership, or repayment obligations.",
  },
  {
    number: 4,
    title: "Lifecycle and liquidity",
    description:
      "Investors access, hold, and trade their positions. Trustees and custodians ensure ongoing compliance and fund protection.",
  },
];

const HowItWorks = () => {
  return (
    <section className="bg-white py-20">
      <div className="page-container">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-24 3xl:gap-32 items-stretch">
          {/* Left: steps content */}
          <div className="flex flex-col">
            <p className="text-primary font-bold tracking-[0.1em] uppercase mb-4 font-matahari text-sm sm:text-base">
              HOW IT WORKS
            </p>
            <h2 className="text-4xl sm:text-5xl 3xl:text-[56px] font-extrabold leading-[1.08] mb-6 text-secondary font-matahari">
              Four steps from asset to investor
            </h2>
            <p className="text-gray-600 mb-10 leading-[160%] text-sm sm:text-base 3xl:text-lg max-w-[520px] 3xl:max-w-[580px]">
              Trovotech provides the structured rails that take a real-world
              productive asset and make it legally investable at scale.
            </p>

            <div className="flex flex-col gap-2 flex-1">
              {STEPS_DATA.map((step, idx) => (
                <div key={idx} className="flex gap-4 sm:gap-6 px-4 py-4">
                  <div className="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 font-semibold text-sm bg-primary text-white shadow-md">
                    {step.number}
                  </div>
                  <div className="flex-1">
                    <h3 className="text-lg sm:text-xl font-bold text-secondary font-sans">
                      {step.title}
                    </h3>
                    <p className="text-sm sm:text-base text-gray-600 leading-[150%] mt-2">
                      {step.description}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Right: image — full height match */}
          <div className="rounded-3xl overflow-hidden min-h-[400px]">
            <img
              src="/img/new-website/howitworks.jpg"
              alt="Investor using Trovo App"
              className="w-full h-full object-cover"
            />
          </div>
        </div>
      </div>
    </section>
  );
};

export default HowItWorks;
