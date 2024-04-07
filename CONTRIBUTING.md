<!-- omit in toc -->
# Contributing to EMF Command Line Interface (CLI)

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them. Please make sure to read the relevant section before making your contribution. It will make it a lot easier for us maintainers and smooth out the experience for all involved. The community looks forward to your contributions. 🎉

> And if you like the project, but just don't have time to contribute, that's fine. There are other easy ways to support the project and show your appreciation, which we would also be very happy about:
> - Star the project
> - Tweet about it
> - Refer this project in your project's readme
> - Mention the project at local meetups and tell your friends/colleagues

<!-- omit in toc -->
## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [I Have a Question](#i-have-a-question)
- [I Want To Contribute](#i-want-to-contribute)
  - [Reporting Bugs](#reporting-bugs)
  - [Suggesting Enhancements](#suggesting-enhancements)
  - [Your First Code Contribution](#your-first-code-contribution)
  - [Improving The Documentation](#improving-the-documentation)
- [Styleguides](#styleguides)
  - [Commit Messages](#commit-messages)


## Code of Conduct

This project and everyone participating in it is governed by the
[EMF CLI Code of Conduct](https://github.com/easy-model-fusion/emf-cliblob/master/CODE_OF_CONDUCT.md).
By participating, you are expected to uphold this code. Please report unacceptable behavior
to [easymodelfusion@gmail.com](mailto:easymodelfusion@gmail.com).


## I Have a Question

> If you want to ask a question, we assume that you have read the available [Documentation](https://easy-model-fusion.github.io/docs/).

Before you ask a question, it is best to search for existing [Issues](https://github.com/easy-model-fusion/emf-cli/issues) that might help you. In case you have found a suitable issue and still need clarification, you can write your question in this issue. It is also advisable to search the internet for answers first.

If you then still feel the need to ask a question and need clarification, we recommend the following:

- Open an [Issue](https://github.com/easy-model-fusion/emf-cli/issues/new).
- Provide as much context as you can about what you're running into.
- Provide project and platform versions (nodejs, npm, etc), depending on what seems relevant.

We will then take care of the issue as soon as possible.

<!--
You might want to create a separate issue tag for questions and include it in this description. People should then tag their issues accordingly.

Depending on how large the project is, you may want to outsource the questioning, e.g. to Stack Overflow or Gitter. You may add additional contact and information possibilities:
- IRC
- Slack
- Gitter
- Stack Overflow tag
- Blog
- FAQ
- Roadmap
- E-Mail List
- Forum
-->

## I Want To Contribute

> ### Legal Notice <!-- omit in toc -->
> When contributing to this project, you must agree that you have authored 100% of the content, that you have the necessary rights to the content and that the content you contribute may be provided under the project license.

### Reporting Bugs

<!-- omit in toc -->
#### Before Submitting a Bug Report

A good bug report shouldn't leave others needing to chase you up for more information. Therefore, we ask you to investigate carefully, collect information and describe the issue in detail in your report. Please complete the following steps in advance to help us fix any potential bug as fast as possible.

- Make sure that you are using the latest version.
- Determine if your bug is really a bug and not an error on your side e.g. using incompatible environment components/versions (Make sure that you have read the [documentation](https://easy-model-fusion.github.io/docs/). If you are looking for support, you might want to check [this section](#i-have-a-question)).
- To see if other users have experienced (and potentially already solved) the same issue you are having, check if there is not already a bug report existing for your bug or error in the [bug tracker](https://github.com/easy-model-fusion/emf-cli/issues?q=label%3Abug).
- Also make sure to search the internet (including Stack Overflow) to see if users outside of the GitHub community have discussed the issue.
- Collect information about the bug:
  - Stack trace (Traceback)
  - OS, Platform and Version (Windows, Linux, macOS, x86, ARM)
  - Version of the interpreter, compiler, SDK, runtime environment, package manager, depending on what seems relevant.
  - Possibly your input and the output
  - Can you reliably reproduce the issue? And can you also reproduce it with older versions?

<!-- omit in toc -->
#### How Do I Submit a Good Bug Report?

> You must never report security related issues, vulnerabilities or bugs including sensitive information to the issue tracker, or elsewhere in public. Instead sensitive bugs must be sent by email to [easymodelfusion@gmail.com](mailto:easymodelfusion@gmail.com).
<!-- You may add a PGP key to allow the messages to be sent encrypted as well. -->

We use GitHub issues to track bugs and errors. If you run into an issue with the project:

- Open an [Issue](https://github.com/easy-model-fusion/emf-cli/issues/new). (Since we can't be sure at this point whether it is a bug or not, we ask you not to talk about a bug yet and not to label the issue.)
- Explain the behavior you would expect and the actual behavior.
- Please provide as much context as possible and describe the *reproduction steps* that someone else can follow to recreate the issue on their own. This usually includes your code. For good bug reports you should isolate the problem and create a reduced test case.
- Provide the information you collected in the previous section.

Once it's filed:

- The project team will label the issue accordingly.
- A team member will try to reproduce the issue with your provided steps. If there are no reproduction steps or no obvious way to reproduce the issue, the team will ask you for those steps and mark the issue as `needs-repro`. Bugs with the `needs-repro` tag will not be addressed until they are reproduced.
- If the team is able to reproduce the issue, it will be marked `needs-fix`, as well as possibly other tags (such as `critical`), and the issue will be left to be [implemented by someone](#your-first-code-contribution).

<!-- You might want to create an issue template for bugs and errors that can be used as a guide and that defines the structure of the information to be included. If you do so, reference it here in the description. -->


### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Easy Model Fusion, **including completely new features and minor improvements to existing functionality**. Following these guidelines will help maintainers and the community to understand your suggestion and find related suggestions.

<!-- omit in toc -->
#### Before Submitting an Enhancement

- Make sure that you are using the latest version.
- Read the [documentation](https://easy-model-fusion.github.io/docs/) carefully and find out if the functionality is already covered, maybe by an individual configuration.
- Perform a [search](https://github.com/easy-model-fusion/emf-cli/issues) to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one.
- Find out whether your idea fits with the scope and aims of the project. It's up to you to make a strong case to convince the project's developers of the merits of this feature. Keep in mind that we want features that will be useful to the majority of our users and not just a small subset. If you're just targeting a minority of users, consider writing an add-on/plugin library.

<!-- omit in toc -->
#### How Do I Submit a Good Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://github.com/easy-model-fusion/emf-cli/issues).

- Use a **clear and descriptive title** for the issue to identify the suggestion.
- Provide a **step-by-step description of the suggested enhancement** in as many details as possible.
- **Describe the current behavior** and **explain which behavior you expected to see instead** and why. At this point you can also tell which alternatives do not work for you.
- You may want to **include screenshots and animated GIFs** which help you demonstrate the steps or point out the part which the suggestion is related to. You can use [this tool](https://www.cockos.com/licecap/) to record GIFs on macOS and Windows, and [this tool](https://github.com/colinkeenan/silentcast) on Linux. <!-- this should only be included if the project has a GUI -->
- **Explain why this enhancement would be useful** to most EMF CLI users. You may also want to point out the other projects that solved it better and which could serve as inspiration.

<!-- You might want to create an issue template for enhancement suggestions that can be used as a guide and that defines the structure of the information to be included. If you do so, reference it here in the description. -->

### Your First Code Contribution

If you're new to contributing to EMF CLI, welcome aboard! Here's a guide to help you get started with your first code contribution:

1. **Set Up Your Development Environment:**
   - Clone the EMF CLI repository to your local machine. 
   - Install any necessary dependencies outlined in the [Documentation](https://easy-model-fusion.github.io/docs/). 
   - Set up your preferred Integrated Development Environment (IDE) for working with the project code. 
2. **Familiarize Yourself with the Codebase:**
   - Take some time to explore the project's directory structure. 
   - Read through existing code to understand how different components interact. 
   - Refer to the [Documentation](https://easy-model-fusion.github.io/docs/) and code comments for guidance.
3. **Find an Issue to Work On:**
   - Browse the project's issue tracker for tasks labeled as `good first issue` or `help wanted` to find beginner-friendly tasks.
   - If you're unsure about what to work on, feel free to ask for guidance in the project's communication channels (e.g., Slack, Discord, or GitHub Discussions).
4. **Make Your Contribution:**
   - Once you've identified an issue to work on, fork the repository and create a new branch for your changes.
   - Implement the necessary changes following the project's coding standards and guidelines (see [Styleguides](#styleguides)).
   - Write tests to ensure that your code functions as expected and doesn't introduce any regressions.
   - Make small, focused commits with clear and descriptive commit messages (see [Commit Messages](#commit-messages)).
   - Once your changes are ready, submit a pull request ([PR](https://github.com/easy-model-fusion/emf-cli/pulls)) to the main repository.
5. **Collaborate and Iterate:**
   - Be open to feedback from maintainers and other contributors.
   - Address any feedback or code review comments promptly and make necessary adjustments to your code.
   - Work iteratively, updating your PR as needed until it meets the project's standards and is ready to be merged.
6. **Celebrate Your Contribution:**
   - Once your PR is merged, celebrate your contribution to the EMF CLI project!

### Improving The Documentation

Improving documentation is a valuable way to contribute to the EMF CLI project. Here are some ways you can help enhance the project's documentation:

1. **Update Existing Documentation:**
    - Review the project's existing [Documentation](https://easy-model-fusion.github.io/docs/) to identify any outdated or inaccurate information.
    - Make corrections, clarifications, or updates to ensure that the documentation reflects the current state of the project.
2. **Add Missing Documentation:**
    - Identify areas of the project that lack sufficient documentation, such as undocumented features, configuration options, or usage examples.
    - Write clear and concise documentation to fill in these gaps and provide valuable information to users and contributors.
3. **Improve Clarity and Readability:**
    - Rewrite sections of the documentation to improve clarity, readability, and organization.
    - Use descriptive headings, bullet points, and examples to make the documentation more accessible to users with varying levels of expertise.
4. **Correct Grammar and Spelling Errors:**
    - Proofread the documentation for grammar, spelling, and punctuation errors.
    - Fix any typos or inconsistencies to maintain the professionalism and credibility of the documentation.
5. **Provide Examples and Tutorials:**
    - Include code [Demos](https://github.com/easy-model-fusion/demos), examples, tutorials, or walkthroughes to help users understand how to use the EMF CLI effectively.
    - Illustrate common use cases or best practices to assist users in achieving their goals with the tool.
6. **Solicit Feedback and Collaboration:**
    - Encourage users and contributors to provide feedback on the documentation.
    - Collaborate with other contributors to review and improve the documentation collaboratively.

Remember to follow the [Documentation](https://easy-model-fusion.github.io/docs/) guidelines and conventions while making contributions to ensure consistency and coherence across the documentation.

## Styleguides
### Commit Messages

When making commits to the EMF CLI repository, follow these guidelines for writing clear and informative commit messages:

1. **Use a Descriptive Title:**
    - Start your commit message with a brief, descriptive title that summarizes the changes introduced by the commit.
2. **Provide Context and Detail:**
    - After the title, provide additional context and detail about the changes in the commit message body.
    - Explain why the changes were necessary and how they address a specific issue or contribute to the project's goals.
3. **Reference Related Issues:**
    - If the commit is related to a specific issue or feature request, reference the corresponding issue or pull request in the commit message.
    - Use keywords like `Fixes`, `Resolves`, or `Addresses` followed by the issue or PR number to automatically close the referenced item upon merging the commit.
4. **Keep Lines Short and Concise:**
    - Limit each line of the commit message to 72 characters or less to ensure readability in various Git tools and interfaces.
    - Use bullet points or paragraphs to organize longer commit messages for improved clarity and structure.
5. **Use Imperative Mood:**
    - Write commit messages in the imperative mood (e.g., "Fix bug" instead of "Fixed bug") to convey commands or actions performed by the commit.
6. **Proofread and Edit:**
    - Before finalizing your commit message, proofread it carefully to ensure accuracy, clarity, and professionalism.
    - Edit the message as needed to remove redundancy, ambiguity, or irrelevant information.

By following these guidelines, you can create meaningful and informative commit messages that facilitate collaboration and maintain clarity throughout the project's version history.

<!-- omit in toc -->
## Attribution
This guide is based on the **contributing-gen**. [Make your own](https://github.com/bttger/contributing-gen)!
