import CTA from "../components/CTA";
import Header from "../components/Navigation";
import Testimonial from "../components/Testimonial";
import Hero from "../components/Hero";
import Footer from "../components/Footer";
import Features from "../components/Features";
import HowItWorks from "../components/HowItWorks";
import { SignedIn, SignedOut, SignInButton, UserButton } from "@clerk/nextjs";

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white">
      <main>
        
        <Hero />
        <Features />
        <HowItWorks />
        <Testimonial />
        <CTA />
      </main>
      <Footer />
    </div>
  );
}
