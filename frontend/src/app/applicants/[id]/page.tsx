"use client";

import { useParams } from "next/navigation";
import { useState, useEffect } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { CalendarIcon, TrophyIcon, Medal, FileText, Mail } from "lucide-react";

interface Applicant {
  _id: string;
  first_name: string;
  last_name: string;
  major: string;
  year: string;
  rating: number;
  ratingCount: number;
  timestamp: string;
  project_id: string;
  wins: number;
  losses: number;
  elo: number;
  resume: any;
  coverLetter: any;
  image: any;
}

export default function ApplicantPage() {
  const params = useParams();
  const applicantId = params?.id as string;

  const [applicant, setApplicant] = useState<Applicant | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchApplicant = async () => {
      try {
        console.log("Fetching applicant:", applicantId);
        const response = await fetch(
          `http://localhost:8080/api/applicants?id=${applicantId}`
        );
        
        if (!response.ok) {
          throw new Error(`Failed to fetch applicant: ${response.statusText}`);
        }
        
        const data = await response.json();
        console.log("Received applicant data:", data);
        setApplicant(data);
      } catch (err: any) {
        console.error("Error fetching applicant:", err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (applicantId) {
      fetchApplicant();
    }
  }, [applicantId]);

  if (loading) return <div className="text-center py-8">Loading applicant details...</div>;
  if (error) return <div className="text-center text-red-500 py-8">{error}</div>;
  if (!applicant) return <div className="text-center py-8">Applicant not found</div>;

  const winRate = applicant.wins + applicant.losses > 0 
    ? ((applicant.wins / (applicant.wins + applicant.losses)) * 100).toFixed(1)
    : "0.0";

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <Card className="mb-8">
        <CardHeader className="space-y-4">
          <div className="flex justify-between items-start">
            <div>
              <CardTitle className="text-3xl">
                {applicant.first_name} {applicant.last_name}
              </CardTitle>
              <div className="text-gray-500 mt-2 space-x-2">
                <Badge variant="secondary">{applicant.major}</Badge>
                <Badge variant="outline">{applicant.year}</Badge>
              </div>
            </div>
          </div>
        </CardHeader>

        <CardContent className="space-y-8">
          {/* Stats Section */}
          <div className="grid grid-cols-3 gap-4">
            <div className="bg-yellow-50 p-4 rounded-lg text-center">
              <TrophyIcon className="h-6 w-6 mx-auto mb-2 text-yellow-500" />
              <div className="text-2xl font-bold">{applicant.elo}</div>
              <div className="text-sm text-gray-600">Elo Rating</div>
            </div>
            <div className="bg-blue-50 p-4 rounded-lg text-center">
              <Medal className="h-6 w-6 mx-auto mb-2 text-blue-500" />
              <div className="text-2xl font-bold">
                {applicant.wins}-{applicant.losses}
              </div>
              <div className="text-sm text-gray-600">Win/Loss Record</div>
            </div>
            <div className="bg-green-50 p-4 rounded-lg text-center">
              <div className="text-2xl font-bold">{winRate}%</div>
              <div className="text-sm text-gray-600">Win Rate</div>
            </div>
          </div>

          <Separator />

          {/* Application Details */}
          <div className="grid md:grid-cols-2 gap-6">
            {/* Resume Section */}
            <div className="bg-gray-50 p-6 rounded-lg">
              <div className="flex items-center mb-4">
                <FileText className="h-5 w-5 mr-2 text-gray-600" />
                <h3 className="text-lg font-semibold">Resume</h3>
              </div>
              <pre className="whitespace-pre-wrap text-sm bg-white p-4 rounded border">
                {JSON.stringify(applicant.resume, null, 2)}
              </pre>
            </div>

            {/* Cover Letter Section */}
            <div className="bg-gray-50 p-6 rounded-lg">
              <div className="flex items-center mb-4">
                <Mail className="h-5 w-5 mr-2 text-gray-600" />
                <h3 className="text-lg font-semibold">Cover Letter</h3>
              </div>
              <pre className="whitespace-pre-wrap text-sm bg-white p-4 rounded border">
                {JSON.stringify(applicant.coverLetter, null, 2)}
              </pre>
            </div>
          </div>

          {/* Additional Info */}
          <div className="bg-gray-50 p-6 rounded-lg">
            <h3 className="text-lg font-semibold mb-4">Additional Information</h3>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <p className="text-gray-600">Rating:</p>
                <p className="font-medium">{applicant.rating}</p>
              </div>
              <div>
                <p className="text-gray-600">Rating Count:</p>
                <p className="font-medium">{applicant.ratingCount}</p>
              </div>
              <div>
                <p className="text-gray-600">Project ID:</p>
                <p className="font-medium">{applicant.project_id}</p>
              </div>
              <div>
                <p className="text-gray-600">Application Date:</p>
                <p className="font-medium">
                  {new Date(applicant.timestamp).toLocaleDateString()}
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
} 