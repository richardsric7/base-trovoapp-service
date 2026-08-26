import PageHeader from "@/component/PageHeader"
import { FaMapMarkerAlt } from "react-icons/fa"
import { FaClock, FaEnvelope, FaPhone } from "react-icons/fa6"
import { FiClock } from "react-icons/fi"


const ContactPage = () => {
  return (
    <>
      <PageHeader text="Contact" />

      <div role="main" className="main">
        <div className="container mx-auto px-4 md:px-12 lg:px-[90px] py-12">

          <div className="flex flex-wrap -mx-4 py-4">
            <div className="w-full lg:w-1/2 px-4 mb-8 lg:mb-0">

              <h2 className="font-bold text-3xl mt-2 mb-0">Contact Us</h2>
              <p className="mb-4">Feel free to ask for details, don&apos;t save any questions!</p>
              
              <form className="contact-form" action="https://formkeep.com/f/9d4b508ef8f9" method="POST">
                <div className="hidden mt-4 bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative">
                  <strong>Success!</strong> Your message has been sent to us.
                </div>

                <div className="hidden mt-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative">
                  <strong>Error!</strong> There was an error sending your message.
                  <span className="text-xs block"></span>
                </div>

                <div className="flex flex-wrap -mx-4 mb-4 mt-6">
                  <div className="w-full lg:w-1/2 px-4 mb-4 lg:mb-0">
                    <label className="block mb-1 text-sm font-medium text-gray-700">Full Name</label>
                    <input type="text" defaultValue="" data-msg-required="Please enter your name." maxLength={100} className="w-full px-3 py-2 border border-gray-300 rounded text-base focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" name="name" required />
                  </div>
                  <div className="w-full lg:w-1/2 px-4">
                    <label className="block mb-1 text-sm font-medium text-gray-700">Email Address</label>
                    <input type="email" defaultValue="" data-msg-required="Please enter your email address." data-msg-email="Please enter a valid email address." maxLength={100} className="w-full px-3 py-2 border border-gray-300 rounded text-base focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" name="email" required />
                  </div>
                </div>
                <div className="flex flex-wrap -mx-4 mb-4">
                  <div className="w-full px-4">
                    <label className="block mb-1 text-sm font-medium text-gray-700">Subject</label>
                    <input type="text" defaultValue="" data-msg-required="Please enter the subject." maxLength={100} className="w-full px-3 py-2 border border-gray-300 rounded text-base focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" name="subject" required />
                  </div>
                </div>
                <div className="flex flex-wrap -mx-4 mb-4">
                  <div className="w-full px-4">
                    <label className="block mb-1 text-sm font-medium text-gray-700">Message</label>
                    <textarea maxLength={5000} data-msg-required="Please enter your message." rows={8} className="w-full px-3 py-2 border border-gray-300 rounded text-base focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500" name="message" required></textarea>
                  </div>
                </div>
                <div className="flex flex-wrap -mx-4 mt-6">
                  <div className="w-full px-4">
                    <input type="submit" value="Send Message" className="bg-primary hover:bg-blue-700 text-white font-bold py-3 px-8 rounded cursor-pointer transition-colors" data-loading-text="Loading..." />
                  </div>
                </div>
              </form>

            </div>
            
            <div className="w-full lg:w-1/2 px-4 lg:pl-12">

              <div className="animate-fade-in" style={{ animationDelay: '800ms' }}>
                <h4 className="text-2xl font-medium mt-2 mb-4">Our <strong className="font-bold text-gray-900">Office</strong></h4>
                <ul className="list-none mt-2 space-y-4">
                  <li className="flex items-center text-gray-700 gap-3">
                    <div className="border p-1 w-8 h-8 shrink-0 rounded-full flex items-center justify-center text-primary"><FaMapMarkerAlt /></div>
                    <span><strong className="text-gray-900 font-semibold">Address:</strong> Lagos, Nigeria</span>
                  </li>
                  <li className="flex items-center text-gray-700 gap-3">
                    <div className="border p-1 w-8 h-8 shrink-0 rounded-full flex items-center justify-center text-primary"><FaPhone /></div>
                    <span><strong className="text-gray-900 font-semibold">Phone:</strong> (+234) 9030009231</span>
                  </li>
                  <li className="flex items-center text-gray-700 gap-3">
                    <div className="border p-1 w-8 h-8 shrink-0 rounded-full flex items-center justify-center text-primary"><FaPhone /></div>
                    <span><strong className="text-gray-900 font-semibold">Phone:</strong> (+234) 8034477988</span>
                  </li>
                  <li className="flex items-center text-gray-700 gap-3">
                    <div className="border p-1 w-8 h-8 shrink-0 rounded-full flex items-center justify-center text-primary"><FaEnvelope /></div>
                    <span><strong className="text-gray-900 font-semibold">Email:</strong> <a href="mailto:info@trovotech.io" className="!text-primary hover:underline">
                        info@trovotech.io</a></span>
                  </li>
                </ul>
              </div>

              <div className="animate-fade-in pt-10" style={{ animationDelay: '950ms' }}>
                <h4 className="text-2xl font-medium mb-4">Business <strong className="font-bold text-gray-900">Hours</strong></h4>
                <ul className="list-none mt-2 space-y-4 text-gray-700">
                  <li className="flex items-center gap-1">
                    <div className="p-1 w-8 h-8 shrink-0 flex items-center justify-center text-black"><FiClock /></div>
                    <span>Monday - Friday - 9am to 5pm</span>
                  </li>
                  <li className="flex items-center gap-1">
                    <div className="p-1 w-8 h-8 shrink-0 flex items-center justify-center text-black"><FiClock /></div>
                    <span>Saturday - 9am to 2pm</span>
                  </li>
                  <li className="flex items-center gap-1">
                    <div className="p-1 w-8 h-8 shrink-0 flex items-center justify-center text-black"><FiClock /></div>
                    <span>Sunday - Closed</span>
                  </li>
                </ul>
              </div>

              {/* <h4 className="text-2xl font-medium pt-10 mb-4">Get in <strong className="font-bold text-gray-900">Touch</strong></h4>
              <p className="text-lg leading-relaxed mb-0 mt-2 text-gray-700">Do you wish to know more about how we can add value to your business? Reach out to us today.</p> */}

            </div>

          </div>

        </div>

      </div>
    </>
  )
}

export default ContactPage