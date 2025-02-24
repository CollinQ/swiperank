"use client";

import { useState } from "react";

export default function NewProject() {
  const [totalApplicants, setTotalApplicants] = useState<number>(0);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      // You'll want to replace this with your actual API endpoint
      const response = await fetch("/api/add", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ totalApplicants }),
      });

      if (!response.ok) {
        throw new Error("Failed to create project");
      }

      // Clear form after successful submission
      setTotalApplicants(0);
      alert("Project created successfully!");
    } catch (error) {
      console.error("Error creating project:", error);
      alert("Failed to create project");
    }
  };

  return (
    <div className="max-w-md mx-auto mt-10 p-6 bg-white rounded-lg shadow-md">
      <h1 className="text-2xl font-bold mb-6">Create New Project</h1>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="totalApplicants"
            className="block text-sm font-medium text-gray-700"
          >
            Total Applicants
          </label>
          <input
            type="number"
            id="totalApplicants"
            value={totalApplicants}
            onChange={(e) => setTotalApplicants(parseInt(e.target.value) || 0)}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            min="0"
            required
          />
        </div>

        <button
          type="submit"
          className="w-full bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
        >
          Create Project
        </button>
      </form>
    </div>
  );
}
