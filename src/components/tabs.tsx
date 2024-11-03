import React from 'react';

type Props = {
  tabList: string[];
  children: React.ReactNode[];
  onTabChanged?: (index: number) => void;
};

export default function Tabs({ tabList, children = [], onTabChanged }: Props) {
  const [openTab, setOpenTab] = React.useState(1);
  const tabItems = tabList.map((item, index) => (
    <li
      key={`${new Date().getTime()}${item.replace(' ', '')}`}
      id={`${new Date().getTime()}${item.replace(' ', '')}`}
      className="-mb-px mr-2 last:mr-0 flex-auto text-center md:text-lg"
    >
      <a
        className={`font-bold px-5 py-3 block leading-normal
        ${openTab === index + 1 ? 'border-b-2 border-primary-800' : ''}`}
        onClick={(e) => {
          e.preventDefault();
          setOpenTab(index + 1);
          onTabChanged && onTabChanged(index + 1);
        }}
        data-toggle="tab"
        href={`#link${index + 1}`}
        role="tablist"
      >
        {item}
      </a>
    </li>
  ));

  const tabBodies = children?.map((body, index) => (
    <div
      className={openTab === index + 1 ? 'block' : 'hidden'}
      key={`${new Date().getTime()}${index}`}
      id={`link${index + 1}`}
    >
      {body}
    </div>
  ));

  return (
    <div className="flex flex-wrap font-matahariRegular w-full">
      <div className="w-full">
        <ul className="flex mb-0 list-none flex-wrap flex-row" role="tablist">
          {tabItems}
        </ul>
        <div className="relative flex flex-col min-w-0 w-full">
          <div className="px-2 py-5 flex-auto">
            <div className="tab-content tab-space">{tabBodies}</div>
          </div>
        </div>
      </div>
    </div>
  );
}

Tabs.defaultProps = {
  children: [],
};
