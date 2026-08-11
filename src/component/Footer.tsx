"use client";

import Image from "next/image";
import Link from "next/link";
import { HiOutlineMail } from "react-icons/hi";
import { IoCallOutline, IoLocationOutline } from "react-icons/io5";
import logo from "@/app/icon.svg";
const Footer = () => {
  return (
    <footer className="bg-dark text-white pt-20 pb-8 mt-auto">
      <div className=" mx-auto px-5 sm:px-8 md:px-12 lg:px-20">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-16 mb-16">
          <div className="footer-about">
            <Link
              href="/"
              aria-label="Trovotech home"
              className="mb-6 inline-block"
            >
              <div className="flex items-center gap-2.5">
                <Image
                  src={logo}
                  alt="trovotech-logo"
                  width={40}
                  height={40}
                  className="h-10 w-10 object-contain"
                />
                <span className="font-matahari-extended text-3xl leading-none text-white">
                  trovotech
                </span>
              </div>
            </Link>
            <div className="social-links flex gap-4">
              <a
                href="https://twitter.com/Trovotechio/"
                className="social-link hover:opacity-80 transition-opacity"
              >
                <img
                  src="/img/new-website/streamline-logos_x-twitter-logo-block.svg"
                  alt="Twitter"
                  className="w-6 h-6"
                />
              </a>
              <a
                href="https://www.linkedin.com/company/trovotech/"
                className="social-link hover:opacity-80 transition-opacity"
              >
                <img
                  src="/img/new-website/pajamas_linkedin.svg"
                  alt="LinkedIn"
                  className="w-6 h-6"
                />
              </a>
              <a
                href="https://www.facebook.com/trovotech/"
                className="social-link hover:opacity-80 transition-opacity"
              >
                <img
                  src="/img/new-website/streamline-plump_facebook-1-solid.svg"
                  alt="Facebook"
                  className="w-6 h-6"
                />
              </a>
              <a
                href="https://www.youtube.com/@trovotechio"
                className="social-link hover:opacity-80 transition-opacity"
              >
                <img
                  src="/img/new-website/streamline_instagram-solid.svg"
                  alt="Instagram"
                  className="w-6 h-6"
                />
              </a>
              <a
                href="https://www.youtube.com/@trovotechio"
                className="social-link hover:opacity-80 transition-opacity"
              >
                <img
                  src="/img/new-website/ri_youtube-fill.svg"
                  alt="YouTube"
                  className="w-6 h-6"
                />
              </a>
            </div>
          </div>

          <div className="footer-links">
            <h4 className="text-lg font-bold mb-6">Quick Links</h4>
            <ul className="space-y-3">
              <li>
                <a
                  href="/terms"
                  className="text-sm text-gray-400 hover:text-white transition-colors"
                >
                  Terms of Use
                </a>
              </li>
              <li>
                <a
                  href="/privacy-policy"
                  className="text-sm text-gray-400 hover:text-white transition-colors"
                >
                  Privacy Policy
                </a>
              </li>
              <li>
                <a
                  href="/support"
                  className="text-sm text-gray-400 hover:text-white transition-colors"
                >
                  Support
                </a>
              </li>
              <li>
                <a
                  href="/faq"
                  className="text-sm text-gray-400 hover:text-white transition-colors"
                >
                  FAQ
                </a>
              </li>
            </ul>
          </div>

          <div className="footer-links">
            <h4 className="text-lg font-bold mb-6">Contact Us</h4>
            <ul className="space-y-4 text-sm text-gray-400">
              <li className="flex items-center gap-1 text-white">
                <IoLocationOutline color="#fff" />
                Lagos, Nigeria
              </li>
              <li className="flex items-center gap-1  text-white">
                <IoCallOutline color="#fff" />
                (+234) 9030009231
              </li>
              <li className="flex items-center gap-1 text-white">
                <HiOutlineMail color="#fff" />
                info@trovotech.io
              </li>
            </ul>
          </div>
          {/* 
          <div className="newsletter">
            <h4 className="text-lg font-bold mb-6">Subscribe Us</h4>
            <form
              className="flex flex-col gap-2"
              onSubmit={(e) => e.preventDefault()}
            >
              <input
                type="email"
                placeholder="Email Address"
                className="w-full px-4 py-3 rounded-lg border border-white/10 bg-white/5 text-white placeholder-gray-400 focus:outline-none focus:border-primary text-sm"
              />
              <button
                type="submit"
                className="w-full bg-primary hover:bg-primary-dark text-white font-medium py-3 rounded-lg cursor-pointer transition-all duration-300 hover:-translate-y-[2px]"
              >
                Subscribe
              </button>
            </form>
          </div> */}
        </div>
      </div>
    </footer>
  );
};

export default Footer;
