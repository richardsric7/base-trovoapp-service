"use client";

import Link from "next/link";
import Image from "next/image";
import { useState, useEffect, useRef } from "react";
import { usePathname } from "next/navigation";
import { RiMenu2Line } from "react-icons/ri";
import logo from "@/app/icon.svg";
import { FaX, FaChevronDown } from "react-icons/fa6";

const navLinks = [
  { label: "Home", href: "/" },
  { label: "About Us", href: "/about" },
  { label: "Ecosystem", href: "/ecosystem" },
];

const docLinks = [
  // {
  //   label: "Product and Services Description",
  //   href: "https://drive.google.com/file/d/1v_BbpTfLR2rebm0HMsnJ55u2bpftMQRb/view",
  // },
  {
    label: "TROV Token Whitepaper",
    href: "https://drive.google.com/file/d/1sGpFZdzusFp3V3GjyYdv797SgwsOvmRg/view",
  },
  {
    label: "Trovo App Whitepaper",
    href: "https://drive.google.com/file/d/16O590h9g8XeHwp9tAjU0hJJ3Z1eYlySk/view",
  },
  // {
  //   label: "Why Build on Bantu Blockchain",
  //   href: "https://drive.google.com/file/d/1Ma92_DQKytJmaFhnUYbjrj5uivFVqAeb/view",
  // },
];

const Navbar = () => {
  const pathname = usePathname();
  const [isScrolled, setIsScrolled] = useState(false);
  const [isOpen, setOpen] = useState(false);
  const [isDocOpen, setIsDocOpen] = useState(false);
  const [isMobileDocOpen, setIsMobileDocOpen] = useState(false);
  const docDropdownRef = useRef<HTMLLIElement>(null);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 50);
    };

    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        docDropdownRef.current &&
        !docDropdownRef.current.contains(event.target as Node)
      ) {
        setIsDocOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const isHome = pathname === "/";
  const isTransparent = isHome && !isScrolled;

  return (
    <>
      <header
        className={`fixed top-0 left-0 w-full z-[1000] transition-all duration-300 ease-in-out ${
          isTransparent
            ? "bg-transparent border-b border-transparent h-20 shadow-none"
            : "bg-white/95 backdrop-blur-md border-b border-gray-200 h-16 shadow-[0_4px_20px_rgba(0,0,0,0.05)]"
        }`}
      >
        <div className="mx-auto px-5 sm:px-8 md:px-12 lg:px-20 flex justify-between items-center h-full">
          <Link
            href="/"
            aria-label="Trovotech home"
            className="border-0 focus:outline-none"
          >
            <div className="flex items-center gap-2.5">
              <Image
                src={logo}
                alt="trovotech-logo"
                width={40}
                height={40}
                priority
                className="h-10 w-10 object-contain"
              />
              <span
                className={`font-matahari-extended text-3xl leading-none transition-colors duration-300 ${
                  isTransparent ? "text-white" : "text-trovo-black"
                }`}
              >
                trovotech
              </span>
            </div>
          </Link>

          {/* Hamburger – mobile only */}
          <RiMenu2Line
            onClick={() => setOpen(true)}
            className={`md:hidden cursor-pointer text-2xl ${isTransparent ? "text-white" : "text-trovo-black"}`}
          />

          {/* Desktop nav */}
          <nav className="hidden md:block h-full">
            <ul className="flex gap-1 items-center h-full">
              {navLinks.map(({ label, href }) => {
                const isActive =
                  href === "/" ? pathname === "/" : pathname.startsWith(href);
                const showLine = isActive && href !== "/";
                return (
                  <li key={href} className="h-full flex items-center relative">
                    <Link
                      href={href}
                      className={`relative h-full flex items-center px-3.5 text-[16px] transition-colors ${
                        isTransparent ? "!text-white" : "!text-trovo-black"
                      } ${isActive ? "font-bold" : "font-normal hover:font-bold"}`}
                    >
                      <span>{label}</span>
                      {showLine && (
                        <span
                          className={`absolute bottom-0 left-0 right-0 h-[3px] rounded-t-sm ${
                            isTransparent ? "bg-white" : "bg-trovo-black"
                          }`}
                        />
                      )}
                    </Link>
                  </li>
                );
              })}

              {/* Documentation Dropdown */}
              <li
                className="relative h-full flex items-center"
                ref={docDropdownRef}
                onMouseEnter={() => setIsDocOpen(true)}
                onMouseLeave={() => setIsDocOpen(false)}
              >
                <button
                  type="button"
                  onClick={() => setIsDocOpen((prev) => !prev)}
                  className={`relative h-full flex items-center gap-1.5 text-[16px] px-3.5 font-normal transition-colors cursor-pointer border-0 bg-transparent hover:font-bold ${
                    isTransparent ? "!text-white" : "!text-trovo-black"
                  }`}
                >
                  <span>Documentation</span>
                  <FaChevronDown
                    className={`text-xs transition-transform duration-200 ${
                      isDocOpen ? "rotate-180" : ""
                    }`}
                  />
                </button>

                {isDocOpen && (
                  <div className="absolute left-0 top-full mt-0 w-64 bg-white rounded-xl shadow-xl border border-gray-100 py-2 z-50 animate-in fade-in slide-in-from-top-2 duration-150">
                    {docLinks.map((doc) => (
                      <a
                        key={doc.label}
                        href={doc.href}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="block px-4 py-2 text-[14px] text-gray-700 hover:text-trovo-black hover:bg-gray-100/60 hover:font-bold transition-all"
                      >
                        {doc.label}
                      </a>
                    ))}
                  </div>
                )}
              </li>

              {/* Contact Us */}
              <li className="h-full flex items-center relative">
                <Link
                  href="/contact"
                  className={`relative h-full flex items-center px-3.5 text-[16px] transition-colors ${
                    isTransparent ? "!text-white" : "!text-trovo-black"
                  } ${
                    pathname.startsWith("/contact")
                      ? "font-bold"
                      : "font-normal hover:font-bold"
                  }`}
                >
                  <span>Contact Us</span>
                  {pathname.startsWith("/contact") && (
                    <span
                      className={`absolute bottom-0 left-0 right-0 h-[3px] rounded-t-sm ${
                        isTransparent ? "bg-white" : "bg-trovo-black"
                      }`}
                    />
                  )}
                </Link>
              </li>
            </ul>
          </nav>
        </div>
      </header>

      {/* Full-screen mobile menu overlay */}
      {isOpen && (
        <div className="md:hidden fixed inset-0 z-[2000] bg-white flex flex-col overflow-y-auto">
          {/* Top bar inside overlay */}
          <div className="flex items-center justify-between px-5 py-5 border-b border-gray-100">
            <Link
              href="/"
              aria-label="Trovotech home"
              onClick={() => setOpen(false)}
            >
              <div className="flex items-center gap-2.5">
                <Image
                  src="/img/favicon.png"
                  alt=""
                  width={40}
                  height={40}
                  className="h-10 w-10 object-contain"
                />
                <span className="font-matahari-extended text-3xl leading-none text-trovo-black">
                  trovotech
                </span>
              </div>
            </Link>
            <FaX
              onClick={() => setOpen(false)}
              className="cursor-pointer text-trovo-black text-xl"
            />
          </div>

          {/* Links */}
          <nav className="flex-1 flex flex-col justify-center px-8 py-6">
            <ul className="flex flex-col gap-6">
              {navLinks.map(({ label, href }) => {
                const isActive =
                  href === "/" ? pathname === "/" : pathname.startsWith(href);
                const showLine = isActive && href !== "/";
                return (
                  <li key={href}>
                    <Link
                      href={href}
                      onClick={() => setOpen(false)}
                      className={`text-[22px] inline-block pb-1 !text-trovo-black hover:font-bold transition-all ${
                        showLine
                          ? "border-b-2 border-trovo-black font-bold"
                          : isActive
                          ? "font-bold"
                          : "font-normal"
                      }`}
                    >
                      {label}
                    </Link>
                  </li>
                );
              })}

              {/* Mobile Documentation Dropdown */}
              <li>
                <button
                  type="button"
                  onClick={() => setIsMobileDocOpen((prev) => !prev)}
                  className="flex items-center justify-between w-full text-[22px] font-normal pb-1 !text-trovo-black hover:font-bold transition-all bg-transparent text-left"
                >
                  <span>Documentation</span>
                  <FaChevronDown
                    className={`text-sm transition-transform duration-200 ${
                      isMobileDocOpen ? "rotate-180" : ""
                    }`}
                  />
                </button>
                {isMobileDocOpen && (
                  <div className="mt-3 pl-4 flex flex-col gap-3 border-l-2 border-primary/20">
                    {docLinks.map((doc) => (
                      <a
                        key={doc.label}
                        href={doc.href}
                        target="_blank"
                        rel="noopener noreferrer"
                        onClick={() => setOpen(false)}
                        className="text-[16px] text-gray-700 hover:text-trovo-black hover:font-bold transition-colors"
                      >
                        {doc.label}
                      </a>
                    ))}
                  </div>
                )}
              </li>

              {/* Contact Us */}
              <li>
                <Link
                  href="/contact"
                  onClick={() => setOpen(false)}
                  className={`text-[22px] inline-block pb-1 !text-trovo-black hover:font-bold transition-all ${
                    pathname.startsWith("/contact")
                      ? "border-b-2 border-trovo-black font-bold"
                      : "font-normal"
                  }`}
                >
                  Contact Us
                </Link>
              </li>
            </ul>
          </nav>
        </div>
      )}
    </>
  );
};

export default Navbar;
