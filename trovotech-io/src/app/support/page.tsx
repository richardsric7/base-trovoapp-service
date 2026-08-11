import PageHeader from '@/component/PageHeader'
import { FaPlane } from 'react-icons/fa6'
import { IoMailOpenOutline, IoPaperPlaneOutline } from 'react-icons/io5'

const SupportPage = () => {
  return (
    <>
      {/* Page Header */}
      <PageHeader text="Get Support" />

      <div role="main" className="mx-[15px]">
        <div className="container mx-auto pt-14">
          <div className="flex flex-wrap justify-center -mx-4 mb-12 pb-3">
            
            <div className="w-full md:w-1/2 lg:w-1/3 px-4 mb-12 lg:mb-0" data-appear-animation="fadeInUpShorter" data-appear-animation-delay="400">
              <div className="flex flex-col text-center  rounded-none bg-gray-50 shadow hover:shadow-lg transition-shadow duration-300">
                <div className="p-6">
                    <IoPaperPlaneOutline className='w-10 h-10 text-primary mx-auto' />

               
                  <h4 className="mt-2 mb-2 text-xl font-bold">Telegram Support</h4>
                  <p className="mb-4">Join our Telegram Support group to report issues and errors you're having with the Wallet.</p>
                  <a href="https://t.me/+9gzT22Q02mswOTRk" target="_blank" className="text-primary font-semibold text-sm">
                  
                  <span className='text-primary'>Join Our Community</span></a>
                </div>
              </div>
            </div>

            <div className="w-full md:w-1/2 lg:w-1/3 px-4 mb-12 lg:mb-0" data-appear-animation="fadeInUpShorter" data-appear-animation-delay="400">
              <div className="flex flex-col text-center rounded-none bg-gray-50 shadow hover:shadow-lg transition-shadow duration-300">
                <div className="p-6">
        <IoMailOpenOutline className='w-10 h-10 text-primary mx-auto'  />

                  <h4 className="mt-2 mb-2 text-xl font-bold">Email</h4>
                  <p className="mb-4">Get in touch with us today via email with any isssues or concerns and we will get back to you.</p>
                  <a href="mailto:support@trovotech.io" className="text-blue-600 font-semibold text-sm"><span className="text-primary">Send an Email</span></a>
                </div>
              </div>
            </div>

          </div>
        </div>
      </div>
    </>
  )
}

export default SupportPage
