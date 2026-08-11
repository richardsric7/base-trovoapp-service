
const TEAM = [
  {
    name: "Thomas Enechi",
    role: "CEO",
    image: "/img/Thomas Enechi https___www.linkedin 1.png",
    linkedin: "https://www.linkedin.com/company/trovotech/",
  },



  {
    name: "Obiefuna Ezenwugo",
    role: "CEO",
    image: "/img/Obiefuna Ezenwugo https___www.linkedin 1.png",
    linkedin: "https://www.linkedin.com/company/trovotech/",
  },
  {
    name: "Babatunde Hassan",
    role: "CEO",
    image: "/img/Babatunde Hassan https___www.linkedin 1.png",
    linkedin: "https://www.linkedin.com/company/trovotech/",
  },
  {
    name: "Ric Richards",
    role: "CEO",
    image: "/img/Ric Richards https___www.linkedin 1.png",
    linkedin: "https://www.linkedin.com/company/trovotech/",
  },
  {
    name: "Ric Richards",
    role: "CEO",
    image: "/img/Ric Richards https___www.linkedin 1.png",
    linkedin: "https://www.linkedin.com/company/trovotech/",
  },

];

const Team = () => {
  return (
    <section className="pt-10  pb-20 bg-white">
      <div className=" mx-auto px-5 sm:px-8 md:px-12 lg:px-20">
        {/* Header */}
        <div className="max-w-[640px] mx-auto mb-16 text-center">
          <p className="text-primary font-bold tracking-[0.1em] uppercase mb-4 font-matahari text-sm sm:text-base">
            MEET THE TEAM
          </p>
          <h2 className="text-4xl sm:text-5xl font-extrabold leading-[1.1] mb-6 text-secondary font-matahari">
            The Team Building This
          </h2>
          <p className="text-gray-600 leading-[160%] text-sm sm:text-base">
            Trovotech is led by experienced professionals across blockchain infrastructure, capital markets, securities law, and regulatory compliance — with deep understanding of African market dynamics and a track record of deploying enterprise technology at global scale.
          </p>
        </div>

        {/* Team grid */}
        <div
          className="flex gap-5 md:gap-8 overflow-x-scroll"
          style={{ scrollbarWidth: "none", msOverflowStyle: "none", WebkitOverflowScrolling: "touch" }}
        >
          {TEAM.map((member, idx) => (
            <div
              key={idx}
              className="group flex-shrink-0 w-[220px] sm:w-[260px] flex flex-col rounded-2xl overflow-hidden"
            >
              {/* Photo */}
              <div className="relative w-full aspect-[4/5] overflow-hidden bg-[#e8ecf4]">
                <img
                  src={member.image}
                  alt={member.name}
                  className="w-full h-full object-cover object-top "
                />
                {/* LinkedIn overlay on hover */}
              </div>
              {/* Info */}
              <div className="p-4 flex flex-col gap-1">
                <h3 className="font-bold text-secondary text-sm sm:text-base leading-snug font-sans">
                  {member.name}
                </h3>
                <p className="text-xs sm:text-sm text-gray-500 leading-snug">
                  {member.role}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default Team;
