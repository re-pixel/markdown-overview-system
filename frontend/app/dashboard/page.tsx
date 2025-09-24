"use client";

import { useState, useEffect } from "react";
import { Upload, FileText, History } from "lucide-react";

interface FileItem {
  name: string;
  uploaded_at: string;
}

export default function DashboardPage() {
  const [file, setFile] = useState<File | null>(null);
  const [model, setModel] = useState("gpt-4");
  const [overview, setOverview] = useState("");
  const [loading, setLoading] = useState(false);
  const [files, setFiles] = useState<FileItem[]>([]);
  const [historyLoading, setHistoryLoading] = useState(false);

  function handleFileChange(e: React.ChangeEvent<HTMLInputElement>) {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  }

  async function handleUpload() {
    if (!file) return;
    setLoading(true);
    setOverview("");

    try {
      const formData = new FormData();
      formData.append("file", file);
      formData.append("model", model);

      const res = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: formData,
        credentials: "include",
      });

      if (res.ok) {
        const data = await res.json();
        setOverview(data.overview || "No overview received.");
        fetchFiles(); // refresh history after upload
      } else {
        setOverview("Error: Failed to generate overview.");
      }
    } catch (err) {
      console.error(err);
      setOverview("Something went wrong while uploading.");
    } finally {
      setLoading(false);
    }
  }

  async function fetchFiles() {
    setHistoryLoading(true);
    try {
      const res = await fetch("http://localhost:8080/files", {
        method: "POST",
        credentials: "include",
      });

      if (res.ok) {
        const data = await res.json();
        setFiles(data.files || []);
      } else {
        console.error("Failed to fetch files history");
      }
    } catch (err) {
      console.error(err);
    } finally {
      setHistoryLoading(false);
    }
  }

  useEffect(() => {
    fetchFiles();
  }, []);

  return (
    <main className="relative flex min-h-screen flex-col items-center bg-gradient-to-br from-blue-50 via-teal-50 to-green-100 p-8 overflow-hidden">
      {/* Animated blobs */}
      <div className="absolute -top-32 -left-32 h-96 w-96 rounded-full bg-teal-300 opacity-30 blur-3xl animate-pulse"></div>
      <div className="absolute -bottom-32 -right-32 h-96 w-96 rounded-full bg-blue-300 opacity-30 blur-3xl animate-pulse delay-1000"></div>

      <section className="relative z-10 grid w-full max-w-6xl gap-8 md:grid-cols-3">
        {/* Main Upload Section */}
        <div className="md:col-span-2 space-y-8">
          <h1 className="text-center text-4xl font-extrabold text-gray-900 drop-shadow-sm">
            File Overview System
          </h1>

          {/* Upload area */}
          <div
            className="flex flex-col items-center justify-center rounded-2xl border-2 border-dashed border-teal-400 bg-white/60 p-10 backdrop-blur-md shadow-lg cursor-pointer hover:border-teal-600 transition"
            onClick={() => document.getElementById("fileInput")?.click()}
          >
            <Upload className="h-12 w-12 text-teal-500 mb-3" />
            <p className="text-gray-700">
              {file ? file.name : "Drag & drop a file here, or click to browse"}
            </p>
            <input
              id="fileInput"
              type="file"
              onChange={handleFileChange}
              className="hidden"
            />
          </div>

          {/* Model selector */}
          <div className="flex flex-col gap-2">
            <label className="font-semibold text-gray-700">Choose Model</label>
            <select
              value={model}
              onChange={(e) => setModel(e.target.value)}
              className="rounded-lg border border-gray-300 p-3 focus:border-teal-500 focus:ring focus:ring-teal-200"
            >
              <option value="gpt-4">GPT-4</option>
              <option value="gpt-3.5">GPT-3.5</option>
              <option value="claude-3">Claude 3</option>
            </select>
          </div>

          {/* Upload button */}
          <button
            onClick={handleUpload}
            disabled={!file || loading}
            className="w-full rounded-lg bg-teal-600 py-3 text-lg font-semibold text-white shadow-md transition hover:bg-teal-700 disabled:opacity-50"
          >
            {loading ? "Processing..." : "Upload & Generate Overview"}
          </button>

          {/* Overview output */}
          <div className="rounded-2xl bg-white/70 backdrop-blur-md p-6 shadow-lg">
            <div className="flex items-center gap-2 mb-3">
              <FileText className="h-6 w-6 text-blue-600" />
              <h2 className="text-xl font-bold text-gray-900">Overview</h2>
            </div>
            <p className="text-gray-700 whitespace-pre-wrap">
              {overview || "Your file overview will appear here."}
            </p>
          </div>
        </div>

        {/* File history sidebar */}
        <aside className="space-y-4 rounded-2xl bg-white/70 backdrop-blur-md p-6 shadow-lg">
          <div className="flex items-center gap-2">
            <History className="h-6 w-6 text-teal-600" />
            <h2 className="text-lg font-bold text-gray-900">File History</h2>
          </div>
          {historyLoading ? (
            <p className="text-gray-600">Loading history...</p>
          ) : files.length === 0 ? (
            <p className="text-gray-600">No files uploaded yet.</p>
          ) : (
            <ul className="max-h-[500px] overflow-y-auto space-y-3">
              {files.map((f, i) => (
                <li
                  key={i}
                  className="cursor-pointer rounded-lg border border-gray-200 bg-white/80 p-3 text-gray-800 shadow-sm transition hover:bg-teal-50"
                  onClick={() => setOverview(`Overview for: ${f.name}\n\n(Click would fetch from backend)`)}
                >
                  <p className="font-medium">{f.name}</p>
                  <p className="text-xs text-gray-500">
                    {new Date(f.uploaded_at).toLocaleString()}
                  </p>
                </li>
              ))}
            </ul>
          )}
        </aside>
      </section>
    </main>
  );
}
