import { notFound } from "next/navigation";

interface Applicant {
  id: string;
  name: string;
  resume: string; // Could be a URL or text excerpt
  coverLetter: string;
  interviewNotes: string;
  customCriteria: {
    [key: string]: string;
  };
}

// Simulated API call to fetch applicant data by ID
async function getApplicant(id: string): Promise<Applicant | null> {
  // In a real implementation, fetch from your API or database.
  // For demonstration, return dummy data.
  const dummyApplicant: Applicant = {
    id: "1",
    name: "Jane Doe",
    resume: "View Resume: https://example.com/resume.pdf",
    coverLetter:
      "Cover Letter: Experienced professional with a background in software development.",
    interviewNotes:
      "Interview Notes: Demonstrated strong problem-solving skills and a team-oriented mindset.",
    customCriteria: {
      "Technical Skills": "Proficient in React, Node.js, and TypeScript",
      "Cultural Fit": "Passionate about innovation and collaboration",
    },
  };

  return dummyApplicant;
}

export default async function ApplicantPage({
  params,
}: {
  params: { id: string };
}) {
  const { id } = await params;
  const applicant = await getApplicant(id);

  if (!applicant) {
    return notFound();
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">
        Applicant Details: {applicant.name}
      </h1>

      <section className="mb-6">
        <h2 className="text-2xl font-semibold mb-2">Resume</h2>
        <p className="text-gray-700">{applicant.resume}</p>
      </section>

      <section className="mb-6">
        <h2 className="text-2xl font-semibold mb-2">Cover Letter</h2>
        <p className="text-gray-700">{applicant.coverLetter}</p>
      </section>

      <section className="mb-6">
        <h2 className="text-2xl font-semibold mb-2">Interview Notes</h2>
        <p className="text-gray-700">{applicant.interviewNotes}</p>
      </section>

      <section>
        <h2 className="text-2xl font-semibold mb-2">Additional Criteria</h2>
        <ul className="list-disc list-inside text-gray-700">
          {Object.entries(applicant.customCriteria).map(([key, value]) => (
            <li key={key}>
              <strong>{key}:</strong> {value}
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
