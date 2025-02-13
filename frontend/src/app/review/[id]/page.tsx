"use client";

import React, { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

interface Applicant {
  id: string;
  name: string;
  resume: string;
  coverLetter: string;
  interviewNotes?: string;
}

// Mock data – replace with your API call as needed
const mockApplicants: Applicant[] = [
  {
    id: "1",
    name: "John Doe",
    resume: "4 years of experience in software development...",
    coverLetter: "I am excited to apply for this position because...",
    interviewNotes: "Strong communication skills",
  },
  {
    id: "2",
    name: "Jane Smith",
    resume: "Recent graduate with a degree in computer science...",
    coverLetter: "I believe your organization aligns with my values...",
    interviewNotes: "Excellent technical background",
  },
];

const CandidatesPage = () => {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [leftApplicant, rightApplicant] = mockApplicants;

  return (
    <div className="flex flex-col items-center justify-center p-4 max-w-7xl mx-auto">
      <div className="flex flex-col md:flex-row gap-4">
        {/* Left Applicant Card */}
        <Card
          onClick={() => setSelectedId(leftApplicant.id)}
          className={`w-full md:w-96 cursor-pointer transition-shadow border-2 ${
            selectedId === leftApplicant.id ? "shadow-xl border-blue-500" : "shadow-sm border-gray-200"
          }`}
        >
          <CardHeader>
            <CardTitle className="text-xl font-bold">{leftApplicant.name}</CardTitle>
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
                    <p className="text-gray-600">{leftApplicant.interviewNotes}</p>
                  </div>
                </>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Right Applicant Card */}
        <Card
          onClick={() => setSelectedId(rightApplicant.id)}
          className={`w-full md:w-96 cursor-pointer transition-shadow border-2 ${
            selectedId === rightApplicant.id ? "shadow-xl border-blue-500" : "shadow-sm border-gray-200"
          }`}
        >
          <CardHeader>
            <CardTitle className="text-xl font-bold">{rightApplicant.name}</CardTitle>
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
                    <p className="text-gray-600">{rightApplicant.interviewNotes}</p>
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
