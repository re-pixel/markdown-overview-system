import Link from "next/link";
import { Upload, LogIn, UserPlus } from "lucide-react";

export default function HomePage() {
  return (
    <main className="relative flex h-screen flex-col items-center justify-center overflow-hidden bg-gradient-to-br from-blue-50 via-teal-50 to-green-100 text-gray-800">
      {/* Animated background blobs */}
      <div className="absolute -top-32 -left-32 h-96 w-96 rounded-full bg-teal-300 opacity-30 blur-3xl animate-pulse"></div>
      <div className="absolute -bottom-32 -right-32 h-96 w-96 rounded-full bg-blue-300 opacity-30 blur-3xl animate-pulse delay-1000"></div>

      {/* Hero Section */}
      <div className="relative text-center max-w-2xl px-6 z-10">
        <h1 className="text-6xl font-extrabold drop-shadow-sm text-gray-900">
          File Overview System
        </h1>
        <p className="mt-4 text-lg text-gray-700">
          Upload your files and let AI generate smart overviews instantly.
        </p>
      </div>

      {/* Call to Action */}
      <div className="relative mt-10 flex gap-6 z-10">
        <Link
          href="/login"
          className="flex items-center gap-2 rounded-full bg-blue-600 px-8 py-3 text-lg font-semibold text-white shadow-md transition hover:scale-105 hover:bg-blue-700"
        >
          <LogIn className="h-5 w-5" />
          Login
        </Link>
        <Link
          href="/register"
          className="flex items-center gap-2 rounded-full bg-green-500 px-8 py-3 text-lg font-semibold text-white shadow-md transition hover:scale-105 hover:bg-green-600"
        >
          <UserPlus className="h-5 w-5" />
          Register
        </Link>
      </div>

      {/* Feature preview */}
      <div className="relative mt-16 grid gap-8 md:grid-cols-2 max-w-4xl px-6 z-10">
        <div className="rounded-2xl bg-white/70 p-6 shadow-lg backdrop-blur-md">
          <Upload className="h-10 w-10 text-blue-500 mb-3" />
          <h2 className="text-2xl font-bold">Upload any file</h2>
          <p className="text-gray-700">
            Drag & drop or browse your files — we handle documents, PDFs, and more.
          </p>
        </div>
        <div className="rounded-2xl bg-white/70 p-6 shadow-lg backdrop-blur-md">
          <h2 className="text-2xl font-bold text-teal-600">AI-powered insights</h2>
          <p className="text-gray-700">
            Get concise summaries and highlights from our LLM instantly.
          </p>
        </div>
      </div>
    </main>
  );
}
