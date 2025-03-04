"use client";

import React, { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

interface Applicant {
  id: string;
  first_name: string;
  last_name: string;
  resume?: string;
  coverLetter?: string;
  interviewNotes?: string;
}

const CandidatesPage = () => {
  const params = useParams();
  const router = useRouter();
  const projectId = params?.id as string;

  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [applicants, setApplicants] = useState<Applicant[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchApplicants = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `http://localhost:8080/api/getTwoForComparison?project_id=${projectId}`
      );
      console.log("Response Status:", response.status);
      if (response.status === 409) {
        // Status code for Conflict
        // Redirect to results page
        router.push(`/results/${projectId}`);
        return;
      }

      if (!response.ok) {
        throw new Error("Failed to fetch applicants");
      }

      const data: Applicant[] = await response.json();
      console.log("API Response Data:", data);

      if (!Array.isArray(data) || data.length < 2) {
        throw new Error("Not enough applicants to compare.");
      }

      setApplicants(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchApplicants();
  }, [projectId]);

  const handleCardSelect = async (winnerId: string, loserId: string) => {
    try {
      // Update Elo ratings
      const response = await fetch("http://localhost:8080/api/updateElo", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          winner_id: winnerId,
          loser_id: loserId,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to update Elo ratings");
      }

      // Clear selection
      setSelectedId(null);

      // Fetch next pair of applicants
      await fetchApplicants();
    } catch (err: any) {
      console.error("Error updating Elo:", err);
      setError(err.message);
    }
  };

  if (loading) {
    return <div className="text-center">Loading applicants...</div>;
  }
  if (error) {
    return <div className="text-center text-red-500">{error}</div>;
  }
  if (applicants.length < 2) {
    return <div className="text-center">Not enough applicants to compare</div>;
  }

  const [leftApplicant, rightApplicant] = applicants;

  return (
    <div className="flex flex-col items-center justify-center p-4 max-w-7xl mx-auto">
      <div className="flex flex-col md:flex-row gap-4">
        {/* Left Applicant Card */}
        <Card
          onClick={() => handleCardSelect(leftApplicant.id, rightApplicant.id)}
          className={`w-full md:w-96 cursor-pointer transition-shadow border-2 ${
            selectedId === leftApplicant.id
              ? "shadow-xl border-blue-500"
              : "shadow-sm border-gray-200"
          }`}
        >
          <CardHeader>
            <CardTitle className="text-xl font-bold">
              {leftApplicant.first_name} {leftApplicant.last_name}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div>
                <p className="font-semibold">Resume</p>
                <p className="text-gray-600">{leftApplicant.resume}</p>
              </div>
              <Separator className="my-2" />
              <div>
                <p className="font-semibold">Cover Letter</p>
                <p className="text-gray-600">{leftApplicant.coverLetter}</p>
              </div>
              {leftApplicant.interviewNotes && (
                <>
                  <Separator className="my-2" />
                  <div>
                    <p className="font-semibold">Interview Notes</p>
                    <p className="text-gray-600">
                      {leftApplicant.interviewNotes}
                    </p>
                  </div>
                </>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Right Applicant Card */}
        <Card
          onClick={() => handleCardSelect(rightApplicant.id, leftApplicant.id)}
          className={`w-full md:w-96 cursor-pointer transition-shadow border-2 ${
            selectedId === rightApplicant.id
              ? "shadow-xl border-blue-500"
              : "shadow-sm border-gray-200"
          }`}
        >
          <CardHeader>
            <CardTitle className="text-xl font-bold">
              {rightApplicant.first_name} {rightApplicant.last_name}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div>
                <p className="font-semibold">Resume</p>
                <p className="text-gray-600">{rightApplicant.resume}</p>
              </div>
              <Separator className="my-2" />
              <div>
                <p className="font-semibold">Cover Letter</p>
                <p className="text-gray-600">{rightApplicant.coverLetter}</p>
              </div>
              {rightApplicant.interviewNotes && (
                <>
                  <Separator className="my-2" />
                  <div>
                    <p className="font-semibold">Interview Notes</p>
                    <p className="text-gray-600">
                      {rightApplicant.interviewNotes}
                    </p>
                  </div>
                </>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default CandidatesPage;
