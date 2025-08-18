import React, { useState } from "react";
import Dropdown from "../../../components/dropdown";
type Step = {
  label: string;
};
const steps: Step[] = [
  { label: "Setup and Compliance" },
  { label: "Asset Information" },
  { label: "Asset Verification Documents" },
  { label: "Asset Token Information" },
];
export function TokenizeAssetForm() {
  const [assetStatus, setAssetStatus] = useState("existing");
  const [offeringType, setOfferingType] = useState("public");
  const [hasDocuments, setHasDocuments] = useState("yes");
  const [confirm, setConfirm] = useState(true);
 const [currentStep, setCurrentStep] = useState(0);
  return (
    <div className="bg-white p-6 w-full mx-auto">
      {/* Title */}
      <h2 className="text-lg md:text-xl font-semibold text-[#0E1F51] mb-6">
        Application to Tokenize Asset
      </h2>
      {/* Steps Progress */}
       <div className="w-full bg-gray-50 rounded-md p-4">
      <div className="flex items-center justify-between">
        {steps.map((step, index) => {
          const isActive = index === currentStep;
          const isCompleted = index < currentStep;

          return (
            <React.Fragment key={step.label}>
              {/* Step Item */}
              <div
                className="flex items-center cursor-pointer"
                onClick={() => setCurrentStep(index)}
              >
                {/* Circle */}
                <div
                  className={`flex items-center justify-center w-6 h-6 rounded-full border 
                  ${
                    isActive
                      ? "bg-blue-500 border-blue-500 text-white"
                      : isCompleted
                      ? "bg-green-500 border-green-500 text-white"
                      : "bg-gray-300 border-gray-300"
                  }`}
                >
                  {isCompleted ? "✓" : ""}
                </div>

                {/* Label */}
                <span
                  className={`ml-2 text-sm font-medium ${
                    isActive ? "text-black font-semibold" : "text-gray-500"
                  }`}
                >
                  {step.label}
                </span>
              </div>

              {/* Line Connector */}
              {index < steps.length - 1 && (
                <div className="flex-1 h-px bg-gray-300 mx-2"></div>
              )}
            </React.Fragment>
          );
        })}
      </div>
    </div>
      {/* <div className="flex items-center justify-between mb-8 bg-[#F2F6F9] py-6 px-4 rounded-2xl">
        {[
          "Setup and Compliance",
          "Asset Information",
          "Asset Verification Documents",
          "Asset Token Information",
        ].map((step, idx) => (
          <div
            key={idx}
            className="flex-1 gap-2 flex justify-center items-center"
          >
            <div
              className={`flex items-center justify-center w-8 h-8 rounded-full text-white text-sm font-bold ${
                idx === 0 ? "bg-[#003A80]" : "bg-gray-300"
              }`}
            >
              {idx + 1}
            </div>
            <span className="text-xs md:text-sm text-center">{step}</span>
          </div>
        ))}
      </div> */}

      {/* Asset Classification */}
      <div className="space-y-6 margin">
        <div className="flex items-center gap-x-6">
          <div className="">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Asset Classification
            </label>
            <p className="text-sm text-gray-500 mb-2 max-w-sm ">
              Please select the classification from the list that suits your
              asset the most
            </p>
          </div>
          {/* <div className="grid md:grid-cols-3 gap-4"> */}
          <div className="flex-1 gap-y-6 flex flex-col">
            <div className="">
               <label className="block font-medium text-[#0E1F51] mb-1">
              Asset Sector
            </label>
              <Dropdown
                label="Asset Sector"
                options={[{ text: "1", value: "2" }]}
                onSelect={() => {}}
              />
            </div>
            <div className="flex align-center flex-1 gap-x-6">
              <div className="flex-1">
                <label className="block font-medium text-[#0E1F51] mb-1">
              Asset Sub-sector
            </label>
                <Dropdown
                  label="Asset Sub-sector"
                  options={[{ text: "1", value: "2" }]}
                  onSelect={() => {}}
                />
              </div>
              <div className="flex-1">
                <label className="block font-medium text-[#0E1F51] mb-1">
              Asset Type
            </label>
                <Dropdown
                  label="Asset Type"
                  options={[{ text: "1", value: "2" }]}
                  onSelect={() => {}}
                />
              </div>
            </div>
          </div>
        </div>
        <div className="flex gap-x-20">
          <div className="">
            <h3 className="block font-medium text-[#0E1F51] mb-1">
              Asset Location
            </h3>
            <p className="block font-medium text-[#0E1F51] mb-1">
              Select Country where the asset is located
            </p>
          </div>

          <div className="flex-1">
            <label className="block font-medium text-[#0E1F51] mb-1">
              Country
            </label>
            <Dropdown
              label="Asset Location"
              options={[{ text: "1", value: "2" }]}
              onSelect={() => {}}
            />
          </div>
        </div>

        <div className="flex gap-x-8">
          <div>
            <label className="block font-medium text-[#0E1F51] mb-1">
              Asset Status
            </label>
            <p className="text-sm text-gray-500 mb-2">
              Select what applies to the current status of your asset
            </p>
          </div>
          <div className="flex-1">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                checked={assetStatus === "existing"}
                onChange={() => setAssetStatus("existing")}
              />
              Asset is already existing
            </label>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                checked={assetStatus === "toBuild"}
                onChange={() => setAssetStatus("toBuild")}
              />
              Asset is yet to be built/acquired
            </label>
          </div>
        </div>

        {/* Offering Type */}
        <div className="flex gap-x-6">
          <div>
            <label className="block font-medium text-[#0E1F51] mb-1">
              Offering Type
            </label>
            <p className="text-sm text-gray-500 max-w-sm mb-2">
              Select the type of offering for your tokenization{" "}
              <span className="italic">
                (private offering is restricted to closed groups)
              </span>
            </p>
          </div>
          <div className="flex-1">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                checked={offeringType === "public"}
                onChange={() => setOfferingType("public")}
              />
              Public
            </label>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                checked={offeringType === "private"}
                onChange={() => setOfferingType("private")}
              />
              Private
            </label>
          </div>
        </div>

        {/* Required Documents */}
        <div className="flex gap-x-6">
          <div>
            <label className="block font-medium text-[#0E1F51] mb-1">
              Required Documents
            </label>
            <p className="text-sm text-gray-500 mb-2">
              Please select the option after checking the requirements
            </p>
          </div>
          <div className=" flex-1">
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                checked={hasDocuments === "yes"}
                onChange={() => setHasDocuments("yes")}
              />
              Yes
            </label>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="radio"
                checked={hasDocuments === "no"}
                onChange={() => setHasDocuments("no")}
              />
              No
            </label>
            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={confirm}
                onChange={() => setConfirm(!confirm)}
              />
              I confirm that tokenizing this asset requires transferring it to a
              licensed Asset Custodian
            </label>
          </div>
        </div>
      </div>

      {/* Buttons */}
      <div className="flex justify-end gap-4 mt-8">
        <button className="px-4 py-2 border rounded-lg hover:bg-gray-100">
          Back
        </button>
        <button className="px-5 py-2 bg-[#003A80] text-white rounded-lg hover:bg-[#00285B]">
          Save and Continue
        </button>
      </div>
    </div>
  );
}
