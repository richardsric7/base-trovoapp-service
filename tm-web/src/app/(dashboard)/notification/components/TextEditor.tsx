"use client";
import React, { useEffect, useRef } from "react";
import Quill from "quill";
import "quill/dist/quill.snow.css";

interface TextEditorProps {
  content: string;
  setContent: React.Dispatch<React.SetStateAction<string>>;
}
const TextEditor: React.FC<TextEditorProps> = ({ content, setContent }) => {
  const toolbarOptions = [
    ["bold", "italic", "underline", "strike"],
    ["blockquote"],
    ["link", "image"],
    [{ list: "ordered" }, { list: "bullet" }, { list: "check" }],
    [{ header: [1, 2, 3, 4, 5, 6, false] }],
    [{ color: [] }, { background: [] }],
    [{ font: [] }],
    [{ align: [] }],
    ["clean"],
  ];

  const contentRef = useRef<HTMLDivElement>(null);
  const quillRef = useRef<Quill | null>(null);

  useEffect(() => {
    if (contentRef.current && !quillRef.current) {
      quillRef.current = new Quill(contentRef.current, {
        theme: "snow",
        modules: {
          toolbar: toolbarOptions,
          clipboard: {
            matchVisual: false,
          },
        } as any,
      });

      quillRef.current.on("text-change", () => {
        setContent(quillRef.current!.root.innerHTML);
      });
    }
  }, [setContent]);

  useEffect(() => {
    if (quillRef.current && content) {
      quillRef.current.root.innerHTML = content;
    }
  }, [content]);

  return (
    <div
      ref={contentRef}
      style={{
        width: "100%",
        minHeight: "360px",
        overflowY: "auto",
      }}
    />
  );
};

export default TextEditor;
