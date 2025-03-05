"use client";

import { useState } from "react";

export default function NewProject() {
  const [name, setName] = useState<string>("");
  const [totalApplicants, setTotalApplicants] = useState<number>(0);
  const [formLink, setFormLink] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      // You'll want to replace this with your actual API endpoint
      const response = await fetch("http://localhost:8080/api/projects", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ name, totalApplicants }),
      });

      if (!response.ok) {
        throw new Error("Failed to create project");
      }
      const data = await response.json();
      setFormLink(data.link);

      // Clear form after successful submission
      setName("");
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
            htmlFor="name"
            className="block text-sm font-medium text-gray-700"
          >
            Project Name
          </label>
          <input
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
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
      {formLink && (
        <div className="mt-4 p-4 bg-gray-50 rounded-md">
          <p className="font-medium">
            Application form has been created! Share the following link with applicants:
          </p>
          <a
            href={formLink}
            target="_blank"
            rel="noopener noreferrer"
            className="text-indigo-600 underline break-words"
          >
            {formLink}
          </a>
        </div>
      )}
    </div>
  );
}
