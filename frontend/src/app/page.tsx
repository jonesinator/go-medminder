"use client";

import { FC, ReactNode, createRef, useState, useEffect } from "react";
import z from "zod";
import { CirclePlusIcon, PencilIcon, RotateCw, XCircleIcon } from "lucide-react";
import { produce } from "immer";

import { apiUrl } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip"

const prescriptionNamesValidator = z.array(z.string());
const prescriptionDataValidator = z.object({
  name: z.string(),
  count: z.number(),
  refill: z.string(),
  rate: z.number(),
  quantity: z.number(),
  updated: z.string(),
});
const prescriptionDataKeys = [...prescriptionDataValidator.keyof().options] as const

const ActionButton: FC<{ tooltip: string, children: ReactNode }> = (props) => (
  <TooltipProvider>
    <Tooltip>
      <TooltipTrigger asChild>
        { props.children }
      </TooltipTrigger>
      <TooltipContent>
        <p>{props.tooltip}</p>
      </TooltipContent>
    </Tooltip>
  </TooltipProvider>
)

export default function page() {
  const [prescriptionData, setPrescriptionData] = useState([] as z.infer<typeof prescriptionDataValidator>[])

  const refresh = () => {
    fetch(apiUrl("rx")).then((response) => {
      response.json().then((responseJson) => {
        const names = prescriptionNamesValidator.parse(responseJson);
        Promise.all(names.map((name) => {
          return fetch(apiUrl(`rx/${name}`)).then((response) => {
            return response.json().then((responseJson => {
              return prescriptionDataValidator.parse(responseJson);
            }))
          })
        })).then((newPrescriptionData) => {
          setPrescriptionData(newPrescriptionData);
        });
      })
    });
  }

  useEffect(() => {
    refresh();
  }, []);

  const createButtonRef = createRef<HTMLButtonElement>();
  const createNameRef = createRef<HTMLInputElement>();
  const createQuantityRef = createRef<HTMLInputElement>();
  const createRateRef = createRef<HTMLInputElement>();

  const handleCreatePrescription = (data: FormData) => {
    if (!createNameRef.current || !createQuantityRef.current || !createRateRef.current || !createButtonRef.current) {
      return;
    }

    createButtonRef.current.disabled = true;

    const object: Record<string, any> = {};
    data.forEach((value, key) => object[key] = value);
    const { name, quantity, rate } = z.object({
      name: z.string(), quantity: z.string(), rate: z.string()
    }).parse(object);
    const body = JSON.stringify({ quantity: parseFloat(quantity), rate: parseFloat(rate) })

    fetch(
      apiUrl(`rx/${name}`),
      { method: "POST", body, headers: { "Content-Type": "application/json" } }
    ).then((response) => {
      if (response.ok) {
        response.json().then((responseJson) => {
          const newPrescription = prescriptionDataValidator.parse(responseJson)
          const newData = produce((draft) => { draft.push(newPrescription) }, prescriptionData);
          setPrescriptionData(newData);
        })
      }
    }).finally(() => {
      createNameRef.current!.value = "";
      createQuantityRef.current!.value = "";
      createRateRef.current!.value = "";
      createButtonRef.current!.disabled = false;
    });
  }

  const handleUpdatePrescription = (data: FormData, name: string) => {
    const object: Record<string, any> = {};
    data.forEach((value, key) => object[key] = value);
    const { quantity, rate } = z.object({ quantity: z.string(), rate: z.string() }).parse(object);
    if (!quantity && !rate) {
      return
    }

    const bodyObject: Record<string, any> = {}
    if (quantity) {
      bodyObject["quantity"] = parseFloat(quantity)
    }
    if (rate) {
      bodyObject["rate"] = parseFloat(rate)
    }
    const body = JSON.stringify(bodyObject)

    fetch(
      apiUrl(`rx/${name}`),
      { method: "PATCH", body, headers: { "Content-Type": "application/json" } }
    ).then((response) => {
      if (response.ok) {
        response.json().then((responseJson) => {
          const newPrescription = prescriptionDataValidator.parse(responseJson)
          const newData = produce((draft) => {
            return draft.map((rx) => {
              if (rx.name !== newPrescription.name) {
                return rx
              } else {
                return newPrescription
              }
            })
          }, prescriptionData);
          setPrescriptionData(newData);
        })
      }
    })
  }

  const handleDeletePrescription = (name: string) => {
    fetch(apiUrl(`rx/${name}`), { method: "DELETE" }).then((response) => {
      if (response.ok) {
        const newData = produce((draft) => {
          return draft.filter((rx) => rx.name !== name)
        }, prescriptionData);
        setPrescriptionData(newData);
      }
    })
  }

  return (<>
    <h1>Go Med-Minder</h1>
    <Table>
      <TableHeader>
        <TableRow>
          {prescriptionDataKeys.map((key) => (<TableHead key={key}>{key}</TableHead>))}
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {prescriptionData.map((rx) => (
          <TableRow key={rx.name}>
            {prescriptionDataKeys.map((key) => <TableCell key={key}>{rx[key]}</TableCell>)}
            <TableCell key="actions" className="flex flex-row">
              <ActionButton tooltip="Edit">
                <Popover>
                  <PopoverTrigger asChild>
                    <PencilIcon className="cursor-pointer mx-1" size={20}/>
                  </PopoverTrigger>
                  <PopoverContent>
                    <form className="flex flex-col" action={(event) => handleUpdatePrescription(event, rx.name)}>
                      <div className="my-2">
                        <label htmlFor="quantity">Quantity:</label>
                        <input ref={createQuantityRef} name="quantity" className="border w-20 ml-4" />
                      </div>
                      <div className="my-2">
                        <label htmlFor="rate">Rate:</label>
                        <input ref={createRateRef} name="rate" className="border w-20 ml-4" />
                      </div>
                      <Button className="m-4" ref={createButtonRef}>
                        Edit
                      </Button>
                    </form>
                  </PopoverContent>
                </Popover>
              </ActionButton>
              <ActionButton tooltip="Delete">
                <XCircleIcon className="text-red-500 cursor-pointer mx-1" size={20} onClick={() => { handleDeletePrescription(rx.name) }}/>
              </ActionButton>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
    <div className="flex flex-row">
      <Popover>
        <PopoverTrigger asChild>
          <Button className="m-4">
            <CirclePlusIcon size={20}/>New Prescription
          </Button>
        </PopoverTrigger>
        <PopoverContent>
          <form className="flex flex-col" action={handleCreatePrescription}>
            <div className="my-2">
              <label htmlFor="name">Name:</label>
              <input ref={createNameRef} name="name" className="border w-20 ml-4" />
            </div>
            <div className="my-2">
              <label htmlFor="quantity">Quantity:</label>
              <input ref={createQuantityRef} name="quantity" className="border w-20 ml-4" />
            </div>
            <div className="my-2">
              <label htmlFor="rate">Rate:</label>
              <input ref={createRateRef} name="rate" className="border w-20 ml-4" />
            </div>
            <Button className="m-4" ref={createButtonRef}>
              <CirclePlusIcon size={20}/>
            </Button>
          </form>
        </PopoverContent>
      </Popover>
      <Button className="m-4" onClick={refresh}>
        <RotateCw size={20}/>Refresh
      </Button>
    </div>
  </>);
}
