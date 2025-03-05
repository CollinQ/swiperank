"use client";

import { useState, useEffect, use } from "react";
import { useSearchParams } from "next/navigation";

export default function ApplyPage() {
    const searchParams = useSearchParams();
    const projectId = searchParams.get("projectId");
    const projectName = searchParams.get("projectName") ? decodeURIComponent(searchParams.get("projectName")!) : "a Project";

    const [name, setName] = useState<string>("");
    const [email, setEmail] = useState<string>("");
    const [resume, setResume] = useState<File | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [success, setSuccess] = useState<boolean>(false);

    useEffect(() => {
        if (!projectId) {
            setError("Invalid project link. No project ID provided")
        }
    }, [projectId]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!projectId || !name || !email || !resume) {
            setError("All fields are required");
            return;
        }

        setLoading(true);
        setError(null);
        setSuccess(false);

        const formData = new FormData();
        formData.append("projectId", projectId);
        formData.append("name", name);
        formData.append("email", email);
        formData.append("resume", resume);

        try {
            const response = await fetch("http://localhost:8080/api/applicants", {
                method: "POST",
                body: formData,
            });

            if (!response.ok) {
                throw new Error("Failed to submit application.");
            }

            setSuccess(true);
            setName("");
            setEmail("");
            setResume(null);
        } catch (error) {
            setError("Failed to submit application. Please try again.");
            console.error("Application error: ", error);
        } finally {
            setLoading(false);
        }
    };

return (
        <div className="max-w-lg mx-auto mt-10 p-6 bg-white rounded-lg shadow-md">
            <h1 className="text-2xl font-bold mb-4">Apply for {projectName}</h1>

            {error && <p className="text-red-500 mb-4">{error}</p>}
            {success && <p className="text-green-500 mb-4">Application submitted successfully!</p>}

            {!error && (
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700">Full Name</label>
                        <input
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700">Email</label>
                        <input
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700">Resume (PDF)</label>
                        <input
                            type="file"
                            accept="application/pdf"
                            onChange={(e) => setResume(e.target.files?.[0] || null)}
                            className="mt-1 block w-full"
                            required
                        />
                    </div>

                    <button
                        type="submit"
                        className="w-full bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                        disabled={loading}
                    >
                        {loading ? "Submitting..." : "Submit Application"}
                    </button>
                </form>
            )}
        </div>
    );
}