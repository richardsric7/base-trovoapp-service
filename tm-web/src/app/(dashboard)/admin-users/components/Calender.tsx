"use client";

import React, { useState } from "react";
import Calendar from "react-calendar";
import "react-calendar/dist/Calendar.css";

type ValuePiece = Date | null;
type Value = ValuePiece | [ValuePiece, ValuePiece];

interface CalenderProps {
  onSelectDate: (date: Date) => void;
}

const Calender: React.FC<CalenderProps> = ({ onSelectDate }) => {
  const [value, onChange] = useState<Value>(new Date());

  const handleDateChange = (newValue: Value) => {
    onChange(newValue);
    if (newValue instanceof Date) {
      onSelectDate(newValue);
    }
  };

  return (
    <div>
      <Calendar onChange={handleDateChange} value={value} />
    </div>
  );
};

export default Calender;
