"use client";

import { useState } from "react";
import Link from "next/link";
import { UserPlus } from "lucide-react";

export default function RegisterPage() {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const res = await fetch("http://localhost:8080/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password }),
      });

      if (res.ok) {
        window.location.href = "/login";
      } else {
        const data = await res.json();
        setError(data.error || "Registration failed");
      }
    } catch (err) {
      console.error(err);
      setError("Something went wrong");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="relative flex h-screen items-center justify-center bg-gradient-to-br from-blue-50 via-teal-50 to-green-100 overflow-hidden">
      {/* Animated blobs */}
      <div className="absolute -top-32 -left-32 h-96 w-96 rounded-full bg-teal-300 opacity-30 blur-3xl animate-pulse"></div>
      <div className="absolute -bottom-32 -right-32 h-96 w-96 rounded-full bg-blue-300 opacity-30 blur-3xl animate-pulse delay-1000"></div>

      <form
        onSubmit={handleSubmit}
        className="relative z-10 w-96 rounded-2xl bg-white/70 backdrop-blur-md p-8 shadow-lg"
      >
        <h1 className="mb-6 text-center text-3xl font-extrabold text-gray-900">
          Register
        </h1>

        {error && <p className="mb-3 text-sm text-red-500">{error}</p>}

        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="mb-4 w-full rounded-lg border border-gray-300 p-3 text-gray-900 font-medium focus:border-blue-500 focus:ring focus:ring-blue-200"
          required
        />

        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="mb-4 w-full rounded-lg border border-gray-300 p-3 text-gray-900 font-medium focus:border-blue-500 focus:ring focus:ring-blue-200"

          required
        />

        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="mb-4 w-full rounded-lg border border-gray-300 p-3 text-gray-900 font-medium focus:border-blue-500 focus:ring focus:ring-blue-200"
          required
        />

        <button
          type="submit"
          disabled={loading}
          className="flex w-full items-center justify-center gap-2 rounded-lg bg-green-600 px-4 py-3 text-white font-semibold shadow-md transition hover:bg-green-700 disabled:opacity-50"
        >
          <UserPlus className="h-5 w-5" />
          {loading ? "Registering..." : "Register"}
        </button>

        <p className="mt-4 text-center text-sm text-gray-600">
          Already have an account?{" "}
          <Link href="/login" className="text-green-600 hover:underline">
            Login here
          </Link>
        </p>
      </form>
    </main>
  );
}
