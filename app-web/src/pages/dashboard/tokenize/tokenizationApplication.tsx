import React, { useMemo } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
type Step = {
  label: string;
};
const steps: Step[] = [
  { label: 'Setup and Compliance' },
  { label: 'Asset Information' },
  { label: 'Asset Verification Documents' },
  { label: 'Asset Token Information' },
];

const getCurrentStepFromPath = (pathname: string): number => {
  if (pathname.includes('/asset-token-information')) {
    return 3;
  }

  if (pathname.includes('/asset-documents')) {
    return 2;
  }

  if (pathname.includes('/asset-information')) {
    return 1;
  }

  if (pathname.includes('/tokenize/apply')) {
    return 0;
  }

  return 0;
};

const getAssetIdFromPath = (pathname: string): string | null => {
  const segments = pathname.split('/').filter(Boolean);
  const applySegmentIndex = segments.indexOf('apply');

  if (applySegmentIndex < 0) {
    return null;
  }

  const valueAfterApply = segments[applySegmentIndex + 1] ?? '';

  if (
    !valueAfterApply ||
    valueAfterApply === 'asset-information' ||
    valueAfterApply === 'asset-documents' ||
    valueAfterApply === 'asset-token-information'
  ) {
    return null;
  }

  return valueAfterApply;
};

export function TokenizationApplication() {
  const location = useLocation();
  const navigate = useNavigate();

  const currentStep = useMemo(() => {
    return getCurrentStepFromPath(location.pathname);
  }, [location.pathname]);

  const assetId = useMemo(() => {
    return getAssetIdFromPath(location.pathname);
  }, [location.pathname]);

  const handleStepClick = (stepIndex: number) => {
    if (stepIndex === 0) {
      navigate(
        assetId
          ? `/dashboard/tokenize/apply/${assetId}`
          : '/dashboard/tokenize/apply',
      );
      return;
    }

    if (!assetId) {
      return;
    }

    const stepPathMap: Record<number, string> = {
      1: `/dashboard/tokenize/apply/${assetId}/asset-information`,
      2: `/dashboard/tokenize/apply/${assetId}/asset-documents`,
      3: `/dashboard/tokenize/apply/${assetId}/asset-token-information`,
    };

    const route = stepPathMap[stepIndex];

    if (route) {
      navigate(route);
    }
  };

  return (
    <div className="bg-white p-6 w-full mx-auto">
      {/* Title */}
      <h2 className="text-lg md:text-xl font-semibold text-[#0E1F51] mb-6">
        Application to Tokenize Asset
      </h2>
      {/* Steps Progress */}
      <div className="w-full bg-gray-50 rounded-md p-4 mb-6">
        <div className="flex items-center justify-between">
          {steps.map((step, index) => {
            const isActive = index === currentStep;
            const isCompleted = index < currentStep;

            return (
              <React.Fragment key={step.label}>
                {/* Step Item */}
                <div
                  className={`flex items-center ${
                    index === 0 || assetId
                      ? 'cursor-pointer'
                      : 'cursor-not-allowed'
                  }`}
                  onClick={() => handleStepClick(index)}
                >
                  {/* Circle */}
                  <div
                    className={`flex items-center justify-center w-6 h-6 rounded-full border 
                  ${
                    isActive
                      ? 'bg-primary-700 border-primary-700 text-white'
                      : isCompleted
                        ? 'bg-green-500 border-green-500 text-white'
                        : 'bg-gray-300 border-gray-300'
                  }`}
                  >
                    {isCompleted ? '✓' : ''}
                  </div>

                  {/* Label */}
                  <span
                    className={`ml-2 text-sm font-medium ${
                      isActive ? 'text-black font-semibold' : 'text-gray-500'
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

      <Outlet />
    </div>
  );
}
