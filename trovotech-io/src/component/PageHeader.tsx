interface HeaderProps{
    text:string;
}

const PageHeader = ({text}:HeaderProps) => {
  return (
      <section className="bg-gradient-to-r from-slate-900 to-slate-800">
        <div className="w-full px-6 pt-[80px] pb-8 md:px-12 md:pt-[110px] md:pb-10 lg:px-[90px]">
          <h1 className="m-0 text-3xl font-bold leading-tight text-white md:text-4xl">
         {text}
          </h1>
        </div>
      </section>
  )
}

export default PageHeader